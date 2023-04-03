package aws

import (
	"context"
	"os"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"

	"github.com/99designs/aws-vault/v7/cli"
	"github.com/99designs/aws-vault/v7/vault"
	"github.com/99designs/keyring"
)

// TryCredentialsFile looks for the access key in the ~/.aws/credentials file.
func TryCredentialsFile() sdk.Importer {
	file := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
	if file == "" {
		file = "~/.aws/credentials"
	}
	return importer.TryFile(file, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
		configPath := os.Getenv("AWS_CONFIG_FILE")
		if configPath != "" {
			if strings.HasPrefix(configPath, "~") {
				configPath = in.FromHomeDir(configPath[1:])
			} else {
				configPath = in.FromRootDir(configPath)
			}
		} else {
			configPath = in.FromHomeDir(".aws", "config") // default config file location
		}
		var configFile *ini.File
		configContent, err := os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			out.AddError(err)
		}
		configFile, err = importer.FileContents(configContent).ToINI()
		if err != nil {
			out.AddError(err)
		}

		for _, section := range credentialsFile.Sections() {
			profileName := section.Name()
			fields := make(map[sdk.FieldName]string)
			if section.HasKey("aws_access_key_id") && section.Key("aws_access_key_id").Value() != "" {
				fields[fieldname.AccessKeyID] = section.Key("aws_access_key_id").Value()
			}

			if section.HasKey("aws_secret_access_key") && section.Key("aws_secret_access_key").Value() != "" {
				fields[fieldname.SecretAccessKey] = section.Key("aws_secret_access_key").Value()
			}

			// read profile configuration from config file
			if configFile != nil {
				configSection := getConfigSectionByProfile(configFile, profileName)
				if configSection != nil {
					if configSection.HasKey("region") && configSection.Key("region").Value() != "" {
						fields[fieldname.DefaultRegion] = configSection.Key("region").Value()
					}
				}
			}

			// add only candidates with required credential fields
			if fields[fieldname.AccessKeyID] != "" && fields[fieldname.SecretAccessKey] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(profileName),
				})
			}
		}
	})
}

// TryAwsVaultCredentials looks for the access key in the user's vaulting backend through AWS Vault.
func TryAwsVaultCredentials() sdk.Importer {
	// Backend types from aws-vault and their respective user-friendly display names
	// Details can be found at https://pkg.go.dev/github.com/99designs/keyring@v1.2.2#section-readme
	backendNames := map[keyring.BackendType]string{
		keyring.SecretServiceBackend: "Secret Service: GNOME Keyring, KWallet",
		keyring.KeychainBackend:      "macOS Keychain",
		keyring.KeyCtlBackend:        "KeyCtl",
		keyring.KWalletBackend:       "KWallet",
		keyring.WinCredBackend:       "Windows Credential Manager",
		keyring.FileBackend:          "Encrypted file",
		keyring.PassBackend:          "Pass",
	}

	// Determine the vaulting backend through AWS_VAULT_BACKEND or from those available on the current OS
	var awsVaultBackend keyring.BackendType
	awsVaultBackendEnvVar := keyring.BackendType(os.Getenv("AWS_VAULT_BACKEND"))
	for i, backendType := range keyring.AvailableBackends() {
		// Default to the first available backend
		if i == 0 {
			awsVaultBackend = backendType
		}
		// If AWS_VAULT_BACKEND matches one of the available backends, use it
		if backendType == awsVaultBackendEnvVar {
			awsVaultBackend = awsVaultBackendEnvVar
			break
		}
	}

	return TryAWSVault(backendNames[awsVaultBackend], func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var keyringConfigDefaults = keyring.Config{
			ServiceName:              "aws-vault",
			LibSecretCollectionName:  "awsvault",
			KWalletAppID:             "aws-vault",
			KWalletFolder:            "aws-vault",
			KeychainTrustApplication: true,
			WinCredPrefix:            "aws-vault",
			KeychainName:             "aws-vault",
			FileDir:                  "~/.awsvault/keys/",
		}

		awsVault := &cli.AwsVault{KeyringConfig: keyringConfigDefaults}
		awsVault.KeyringBackend = string(awsVaultBackend)
		keyring, err := awsVault.Keyring()
		if err != nil {
			out.AddError(err)
			return
		}

		// Use the CredentialKeyring struct from aws-vault to retrieve vaulting backend credentials
		credentialKeyring := &vault.CredentialKeyring{Keyring: keyring}

		// Load the AWS config file (default location at ~/.aws/config)
		awsConfigFile, err := awsVault.AwsConfigFile()
		if err != nil {
			out.AddError(err)
			return
		}
		configLoader := &vault.ConfigLoader{File: awsConfigFile}

		// Get the region specified for the "default" profile
		var defaultRegion string
		if defaultSection, ok := awsConfigFile.ProfileSection("default"); ok {
			defaultRegion = defaultSection.Region
		}

		// Iterate through the profiles in the AWS config file and
		// import any matching credentials stored in the vaulting backend
		for _, profileName := range awsConfigFile.ProfileNames() {
			profileFound, _ := credentialKeyring.Has(profileName)
			if !profileFound {
				continue
			}

			creds, err := credentialKeyring.Get(profileName)
			if err != nil {
				continue
			}

			profileConfig, err := configLoader.GetProfileConfig(profileName)
			if err != nil {
				continue
			}

			fields := make(map[sdk.FieldName]string)
			fields[fieldname.AccessKeyID] = creds.AccessKeyID
			fields[fieldname.SecretAccessKey] = creds.SecretAccessKey

			// If a region is specified for the AWS profile, use it.
			// Otherwise, use the "default" profile region if it's specified
			if profileConfig.Region != "" {
				fields[fieldname.DefaultRegion] = profileConfig.Region
			} else if defaultRegion != "" {
				fields[fieldname.DefaultRegion] = defaultRegion
			}

			// Only add candidates with required credential fields
			if fields[fieldname.AccessKeyID] != "" && fields[fieldname.SecretAccessKey] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(profileName),
				})
			}
		}
	})
}

func TryAWSVault(keyringBackend string, result func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportAttempt)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		if keyringBackend == "" {
			return
		}

		attempt := out.NewAttempt(importer.SourceOther(keyringBackend, ""))

		result(ctx, in, attempt)
	}
}
