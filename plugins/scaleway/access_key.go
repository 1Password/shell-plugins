package scaleway

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://www.scaleway.com/en/docs/iam/how-to/create-api-keys/"),
		ManagementURL: sdk.URL("https://console.scaleway.com/iam/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccessKeyID,
				MarkdownDescription: "Access key ID",
				Secret:              false,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Length: 20,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.SecretAccessKey,
				MarkdownDescription: "Secret key.",
				Secret:              true,
				Optional:            false,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "Default region.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.DefaultZone,
				MarkdownDescription: "Default zone.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "Project ID.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.OrgID,
				MarkdownDescription: "Organization ID.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
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
	"SCW_DEFAULT_REGION":          fieldname.DefaultRegion,
	"SCW_DEFAULT_ZONE":            fieldname.DefaultZone,
	"SCW_DEFAULT_PROJECT_ID":      fieldname.ProjectID,
	"SCW_DEFAULT_ORGANIZATION_ID": fieldname.OrgID,
}

func TryScalewayConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/scw/config.yaml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AccessKey == "" || config.SecretKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:     config.AccessKey,
				fieldname.SecretAccessKey: config.SecretKey,
				fieldname.DefaultRegion:   config.DefaultRegion,
				fieldname.DefaultZone:     config.DefaultZone,
				fieldname.ProjectID:       config.ProjectID,
				fieldname.OrgID:           config.OrganizationID,
			},
		})
	})
}

type Config struct {
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	DefaultRegion  string `yaml:"default_region,omitempty"`
	DefaultZone    string `yaml:"default_zone,omitempty"`
	ProjectID      string `yaml:"default_project_id,omitempty"`
	OrganizationID string `yaml:"default_organization_id,omitempty"`
}
