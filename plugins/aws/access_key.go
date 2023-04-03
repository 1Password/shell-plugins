package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
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
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "The default region to use for this access key.",
				Optional:            true,
			},
			{
				Name:                fieldname.OneTimePassword,
				MarkdownDescription: "The one-time code value for MFA authentication.",
				Optional:            true,
			},
			{
				Name:                fieldname.MFASerial,
				MarkdownDescription: "ARN of the MFA serial number to use to generate temporary STS credentials if the item contains a TOTP setup.",
				Optional:            true,
			},
		},
		DefaultProvisioner: AWSProvisioner(),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"AMAZON_ACCESS_KEY_ID":     fieldname.AccessKeyID,
				"AMAZON_SECRET_ACCESS_KEY": fieldname.SecretAccessKey,
				"AWS_DEFAULT_REGION":       fieldname.DefaultRegion,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"AWS_ACCESS_KEY":     fieldname.AccessKeyID,
				"AWS_SECRET_KEY":     fieldname.SecretAccessKey,
				"AWS_DEFAULT_REGION": fieldname.DefaultRegion,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"AWS_ACCESS_KEY":     fieldname.AccessKeyID,
				"AWS_ACCESS_SECRET":  fieldname.SecretAccessKey,
				"AWS_DEFAULT_REGION": fieldname.DefaultRegion,
			}),
			TryCredentialsFile(),
			TryAWSVaultBackends(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"AWS_ACCESS_KEY_ID":     fieldname.AccessKeyID,
	"AWS_SECRET_ACCESS_KEY": fieldname.SecretAccessKey,
	"AWS_DEFAULT_REGION":    fieldname.DefaultRegion,
}
