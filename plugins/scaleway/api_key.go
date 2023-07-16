package scaleway

import (
	"context"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys"),
		ManagementURL: sdk.URL("https://console.scaleway.com/iam/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessKeyID,
				MarkdownDescription: "The ID of the API Key used to authenticate to Scaleway.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 20,
					Prefix: "SCW",
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.SecretAccessKey,
				MarkdownDescription: "The secret access key used to authenticate to Scaleway.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.DefaultOrganization,
				MarkdownDescription: "The default organization ID to use for this access key.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "The default region to use for this access key.",
				Optional:            true,
			},
			{
				Name:                fieldname.DefaultZone,
				MarkdownDescription: "The default zone to use for this access key.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryScalewayConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SCW_ACCESS_KEY":              fieldname.AccessKeyID,
	"SCW_SECRET_KEY":              fieldname.SecretAccessKey,
	"SCW_DEFAULT_ORGANIZATION_ID": fieldname.DefaultOrganization,
	"SCW_DEFAULT_REGION":          fieldname.DefaultRegion,
	"SCW_DEFAULT_ZONE":            fieldname.DefaultZone,
}

func TryScalewayConfigFile() sdk.Importer {
	file := os.Getenv("SCW_CONFIG_PATH")
	if file == "" {
		file = "~/.config/scw/config.yaml"
	}
	return importer.TryFile(file, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AccessKey == "" || config.SecretKey == "" {
			return
		}

		fields := make(map[sdk.FieldName]string)
		fields[fieldname.AccessKeyID] = config.AccessKey
		fields[fieldname.SecretAccessKey] = config.SecretKey
		if config.DefaultOrganizationID != "" {
			fields[fieldname.DefaultOrganization] = config.DefaultOrganizationID
		}
		if config.DefaultRegion != "" {
			fields[fieldname.DefaultRegion] = config.DefaultRegion
		}
		if config.DefaultZone != "" {
			fields[fieldname.DefaultZone] = config.DefaultZone
		}
		out.AddCandidate(sdk.ImportCandidate{
			Fields:   fields,
			NameHint: importer.SanitizeNameHint("default"),
		})

		for profileName, profile := range config.Profiles {
			profileFields := make(map[sdk.FieldName]string)
			profileFields[fieldname.AccessKeyID] = profile.AccessKey
			profileFields[fieldname.SecretAccessKey] = profile.SecretKey
			if profile.DefaultOrganizationID != "" {
				profileFields[fieldname.DefaultOrganization] = profile.DefaultOrganizationID
			}
			if profile.DefaultRegion != "" {
				profileFields[fieldname.DefaultRegion] = profile.DefaultRegion
			}
			if profile.DefaultZone != "" {
				profileFields[fieldname.DefaultZone] = profile.DefaultZone
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields:   profileFields,
				NameHint: importer.SanitizeNameHint(profileName),
			})
		}
	})
}

type Config struct {
	AccessKey             string             `yaml:"access_key"`
	SecretKey             string             `yaml:"secret_key"`
	DefaultOrganizationID string             `yaml:"default_organization_id"`
	DefaultRegion         string             `yaml:"default_region"`
	DefaultZone           string             `yaml:"default_zone"`
	Profiles              map[string]Profile `yaml:"profiles"`
}

type Profile struct {
	AccessKey             string `yaml:"access_key"`
	SecretKey             string `yaml:"secret_key"`
	DefaultOrganizationID string `yaml:"default_organization_id"`
	DefaultRegion         string `yaml:"default_region"`
	DefaultZone           string `yaml:"default_zone"`
}
