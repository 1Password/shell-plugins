package aws

import (
	"context"
	"os"

	"github.com/99designs/aws-vault/v7/cli"
	"github.com/99designs/aws-vault/v7/vault"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TryAwsVaultCredentials() sdk.Importer {
	// Read config file from the location set in AWS_CONFIG_FILE env var or from ~/.aws/config
	// config, _ := vault.LoadConfigFromEnv()
	configPath := os.Getenv("AWS_CONFIG_FILE")
	if configPath == "" {
		configPath = "~/.aws/config"
	}

	return importer.TryVault(configPath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// Determine the vaulting backend through AWS_VAULT_BACKEND or based on the OS
		awsVault := &cli.AwsVault{}
		if awsVaultBackend := os.Getenv("AWS_VAULT_BACKEND"); awsVaultBackend != "" {
			awsVault.KeyringBackend = awsVaultBackend
		}
		keyring, err := awsVault.Keyring()
		if err != nil {
			out.AddError(err)
			return
		}

		// Use the CredentialKeyring struct from aws-vault to retrieve vaulting backend credentials
		credentialKeyring := &vault.CredentialKeyring{Keyring: keyring}

		configFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		// Get the region specified for the "default" profile
		var defaultRegion string
		for _, section := range configFile.Sections() {
			if section.Name() == "default" {
				if section.HasKey("region") && section.Key("region").Value() != "" {
					defaultRegion = section.Key("region").Value()
					break
				}
			}
		}

		// Iterate through the profiles in the AWS config file and
		// import any matching credentials stored in the vaulting backend
		for _, section := range configFile.Sections() {
			profileName := getConfigSectionProfileName(section.Name())
			if profileName != "default" {
				creds, err := credentialKeyring.Get(profileName)
				if err == nil {
					fields := make(map[sdk.FieldName]string)
					fields[fieldname.AccessKeyID] = creds.AccessKeyID
					fields[fieldname.SecretAccessKey] = creds.SecretAccessKey

					if section.HasKey("region") && section.Key("region").Value() != "" {
						fields[fieldname.DefaultRegion] = section.Key("region").Value()
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
			}
		}
	})
}
