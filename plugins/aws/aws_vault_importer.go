package aws

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"

	"github.com/99designs/aws-vault/v7/cli"
	"github.com/99designs/aws-vault/v7/vault"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TryAwsVaultCredentials() sdk.Importer {
	// Read config file from the location set in AWS_CONFIG_FILE env var or from  ~/.aws/config
	configPath := os.Getenv("AWS_CONFIG_FILE")
	if configPath == "" {
		home, _ := homedir.Dir()
		configPath = filepath.Join(home, ".aws", "config") // default config file location
	}

	return importer.TryFile(configPath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// Determine whether aws-vault is being used by the user by
		// checking if the /usr/local/bin/aws-vault executable exists
		if _, err := os.Stat(filepath.Join("usr", "local", "bin", "aws-vault")); err != nil {
			return
		}

		// Determine the vaulting backend through AWS_VAULT_BACKEND or the OS
		awsVault := &cli.AwsVault{}
		if awsVaultBackend := os.Getenv("AWS_VAULT_BACKEND"); awsVaultBackend != "" {
			awsVault.KeyringBackend = awsVaultBackend
		}
		keyring, err := awsVault.Keyring()
		if err != nil {
			out.AddError(err)
			return
		}

		// Use the aws-vault CredentialKeyring struct to retrieve vaulting backend credentials
		credentialKeyring := &vault.CredentialKeyring{Keyring: keyring}

		var configFile *ini.File
		configContent, err := os.ReadFile(configPath)
		if err != nil && !os.IsNotExist(err) {
			out.AddError(err)
			return
		}
		configFile, err = importer.FileContents(configContent).ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		// Iterate through the profiles in the AWS config file and
		// import any matching credentials stored in the vaulting backend
		for _, section := range configFile.Sections() {
			profileName := section.Name()
			if strings.HasPrefix(profileName, "profile ") {
				profileName = strings.TrimPrefix(profileName, "profile ")
			} else {
				continue
			}

			creds, err := credentialKeyring.Get(profileName)
			if err == nil {
				fields := make(map[sdk.FieldName]string)
				fields[fieldname.AccessKeyID] = creds.AccessKeyID
				fields[fieldname.SecretAccessKey] = creds.SecretAccessKey

				if section.HasKey("region") && section.Key("region").Value() != "" {
					fields[fieldname.DefaultRegion] = section.Key("region").Value()
				}

				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(profileName),
				})
			}
		}
	})
}
