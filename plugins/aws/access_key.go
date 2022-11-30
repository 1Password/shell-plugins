package aws

import (
	"context"
	"gopkg.in/ini.v1"
	"os"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	FieldNameDefaultRegion = "Default Region"
	FieldNameSerialNumber  = "MFA Serial"
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
				Name:                FieldNameSerialNumber,
				MarkdownDescription: "The identification number of the MFA device that is associated with the user who is making the GetSessionToken call, usually an Amazon Resource Name (ARN) for a virtual device (such as arn:aws:iam::123456789012:mfa/user).",
				Optional:            true,
			},
		},
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
		if err != nil {
			out.AddError(err)
			return
		}
		if configContent != nil {
			configFile, err = importer.FileContents(configContent).ToINI()
			if err != nil {
				out.AddError(err)
				return
			}
		}

		for _, section := range credentialsFile.Sections() {
			profileName := section.Name()
			fields := make(map[string]string)
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
						fields[FieldNameDefaultRegion] = configSection.Key("region").Value()
					}
				}
			}

			// add only candidates with required credential fields
			if fields[fieldname.AccessKeyID] != "" && fields[fieldname.SecretAccessKey] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: profileName,
				})
			}
		}
	})
}
