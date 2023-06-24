package aws

import (
	"context"
	"os"
	"path/filepath"
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

// Backend types from AWS Vault and their respective user-friendly display names
// Details can be found at https://pkg.go.dev/github.com/99designs/keyring@v1.2.2#section-readme
var backendNames = map[keyring.BackendType]string{
	keyring.SecretServiceBackend: "Secret Service: GNOME Keyring, KWallet",
	keyring.KeychainBackend:      "macOS Keychain",
	keyring.KeyCtlBackend:        "KeyCtl",
	keyring.KWalletBackend:       "KWallet",
	keyring.WinCredBackend:       "Windows Credential Manager",
	keyring.FileBackend:          "Encrypted file",
	keyring.PassBackend:          "Pass",
}

// Default keyring config values, based on the default values used by AWS Vault
// Refer to the aws-vault codebase for more context: https://github.com/99designs/aws-vault/blob/master/cli/global.go
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

// TryAWSVaultBackends looks for AWS credentials in the user's vaulting backends, using functionality provided by AWS Vault.
func TryAWSVaultBackends() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		// Retrieve all available vaulting backends on the current OS
		availableBackends := keyring.AvailableBackends()
		if len(availableBackends) == 0 {
			return
		}

		// Search through each available vaulting backend for AWS credentials
		for _, backendType := range availableBackends {
			attempt := out.NewAttempt(importer.SourceOther(backendNames[backendType], ""))

			awsVault := &cli.AwsVault{KeyringConfig: keyringConfigDefaults}
			awsVault.KeyringBackend = string(backendType)

			keyring, err := awsVault.Keyring()
			if err != nil {
				attempt.AddError(err)
				return
			}

			// Use the CredentialKeyring struct from aws-vault to retrieve vaulting backend credentials
			credentialKeyring := &vault.CredentialKeyring{Keyring: keyring}

			profilesInfo, err := GetProfilesInfo()
			if err != nil {
				attempt.AddError(err)
				return
			}

			// Iterate through the profiles in the AWS config file and
			// import any matching credentials stored in the vaulting backend
			for _, profile := range profilesInfo {
				profileFound, err := credentialKeyring.Has(profile.Name)
				if err != nil {
					attempt.AddError(err)
					continue
				}
				if !profileFound {
					continue
				}

				creds, err := credentialKeyring.Get(profile.Name)
				if err != nil {
					attempt.AddError(err)
					continue
				}

				fields := make(map[sdk.FieldName]string)
				fields[fieldname.AccessKeyID] = creds.AccessKeyID
				fields[fieldname.SecretAccessKey] = creds.SecretAccessKey
				fields[fieldname.MFASerial] = profile.MfaSerial
				fields[fieldname.Region] = profile.Region
				// Only add candidates with required credential fields
				if fields[fieldname.AccessKeyID] != "" && fields[fieldname.SecretAccessKey] != "" {
					attempt.AddCandidate(sdk.ImportCandidate{
						Fields:   fields,
						NameHint: importer.SanitizeNameHint(profile.Name),
					})
				}
			}
		}
	}
}

type ProfileInfoToImport struct {
	Name      string
	MfaSerial string
	Region    string
}

func GetProfilesInfo() ([]ProfileInfoToImport, error) {
	file := os.Getenv("AWS_CONFIG_FILE")
	if file == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		file = filepath.Join(home, "/.aws/config")
	}

	f, err := ini.LoadSources(ini.LoadOptions{
		AllowNestedValues:   true,
		InsensitiveSections: false,
		InsensitiveKeys:     true,
	}, file)
	if err != nil {
		return nil, err
	}

	const (
		configFileRegionKey    = "region"
		profileNamePrefix      = "profile "
		configFileMfaSerialKey = "mfa_serial"
	)

	// Get the region specified for the "default" profile
	var defaultRegion string
	if f.HasSection(defaultProfileName) {
		if defaultSection, err := f.GetSection(defaultProfileName); err != nil && defaultSection.HasKey(configFileRegionKey) {
			key, err := defaultSection.GetKey(configFileRegionKey)
			if err != nil {
				return nil, err
			}
			defaultRegion = key.String()
		}
	}

	var profiles []ProfileInfoToImport
	for _, section := range f.Sections() {
		var region, mfaSerial string
		if section.HasKey(configFileMfaSerialKey) {
			key, err := section.GetKey(configFileMfaSerialKey)
			if err != nil {
				return nil, err
			}
			mfaSerial = key.String()
		}
		if section.HasKey(configFileRegionKey) {
			key, err := section.GetKey(configFileRegionKey)
			if err != nil {
				return nil, err
			}
			region = key.String()
		} else {
			region = defaultRegion
		}
		profiles = append(profiles, ProfileInfoToImport{
			Name:      strings.TrimPrefix(section.Name(), profileNamePrefix),
			MfaSerial: mfaSerial,
			Region:    region,
		})
	}

	return profiles, nil
}
