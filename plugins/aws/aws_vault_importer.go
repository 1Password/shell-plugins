package aws

import (
	"context"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/99designs/aws-vault/v7/cli"
	"github.com/99designs/aws-vault/v7/vault"
)

// TryAwsVaultCredentials looks for the access key in the user's vaulting backend through AWS Vault.
func TryAwsVaultCredentials() sdk.Importer {
	return TryVault(func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportAttempt) {
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
			if profileName != "default" {
				creds, err := credentialKeyring.Get(profileName)
				if err == nil {
					fields := make(map[sdk.FieldName]string)
					fields[fieldname.AccessKeyID] = creds.AccessKeyID
					fields[fieldname.SecretAccessKey] = creds.SecretAccessKey

					profileConfig, err := configLoader.GetProfileConfig(profileName)
					if err != nil {
						out.AddError(err)
						return
					}

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
			}
		}
	})
}

func TryVault(result func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportAttempt)) sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		attempt := out.NewAttempt(sdk.ImportSource{AWSVault: true})
		result(ctx, in, attempt)
	}
}
