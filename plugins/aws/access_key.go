package aws

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	FieldNameDefaultRegion = "Default Region"
	FieldNameMFASerial     = "MFA Serial"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html"),
		ManagementURL: sdk.URL("https://console.aws.amazon.com/iam"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessKeyID,
				MarkdownDescription: "The ID of the access key used to authenticate to AWS.",
				Composition: &schema.ValueComposition{
					Length: 20,
					Prefix: "AKIA",
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.SecretAccessKey,
				MarkdownDescription: "The secret access key used to authenticate to AWS.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                FieldNameDefaultRegion,
				MarkdownDescription: "The default region to use for this access key.",
				Optional:            true,
			},
			{
				Name:                fieldname.OneTimePassword,
				MarkdownDescription: "The one-time code value for MFA authentication.",
				Optional:            true,
			},
			{
				Name:                FieldNameMFASerial,
				MarkdownDescription: "ARN of the MFA serial number to use to generate temporary STS credentials if the item contains a TOTP setup.",
				Optional:            true,
			},
		},
		CacheKeys:   []string{"sts"},
		Provisioner: AWSProvisioner(),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(officialEnvVarMapping),
			importer.TryEnvVarPair(map[string]string{
				fieldname.AccessKeyID:     "AMAZON_ACCESS_KEY_ID",
				fieldname.SecretAccessKey: "AMAZON_SECRET_ACCESS_KEY",
				FieldNameDefaultRegion:    "AWS_DEFAULT_REGION",
			}),
			importer.TryEnvVarPair(map[string]string{
				fieldname.AccessKeyID:     "AWS_ACCESS_KEY",
				fieldname.SecretAccessKey: "AWS_SECRET_KEY",
				FieldNameDefaultRegion:    "AWS_DEFAULT_REGION",
			}),
			importer.TryEnvVarPair(map[string]string{
				fieldname.AccessKeyID:     "AWS_ACCESS_KEY",
				fieldname.SecretAccessKey: "AWS_ACCESS_SECRET",
				FieldNameDefaultRegion:    "AWS_DEFAULT_REGION",
			}),
			TryCredentialsFile(),
		),
	}
}

var officialEnvVarMapping = map[string]string{
	fieldname.AccessKeyID:     "AWS_ACCESS_KEY_ID",
	fieldname.SecretAccessKey: "AWS_SECRET_ACCESS_KEY",
	FieldNameDefaultRegion:    "AWS_DEFAULT_REGION",
}

// TryCredentialsFile looks for the access key in the ~/.aws/credentials file.
func TryCredentialsFile() sdk.Importer {
	return importer.TryFile("~/.aws/credentials", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		var profiles []string
		for _, v := range credentialsFile.Sections() {
			if len(v.Keys()) != 0 {
				profiles = append(profiles, v.Name())
			}
		}

		for _, profile := range profiles {
			cfg, err := config.LoadSharedConfigProfile(ctx, profile)
			if err != nil {
				out.AddError(err)
				return
			}

			if cfg.Credentials.AccessKeyID == "" || cfg.Credentials.SecretAccessKey == "" {
				continue
			}

			fields := map[string]string{
				fieldname.AccessKeyID:     cfg.Credentials.AccessKeyID,
				fieldname.SecretAccessKey: cfg.Credentials.SecretAccessKey,
			}

			if cfg.Region != "" {
				fields[FieldNameDefaultRegion] = cfg.Region
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields:   fields,
				NameHint: profile,
			})
		}
	})
}
