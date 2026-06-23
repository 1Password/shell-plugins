package descope

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ManagementKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.ManagementKey,
		DocsURL:       sdk.URL("https://docs.descope.com/cli/descope"),
		ManagementURL: sdk.URL("https://app.descope.com/settings/company/managementkeys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "Project ID for the current Descope Project.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.ManagementKey,
				MarkdownDescription: "Management Key used to authenticate to Descope.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DESCOPE_PROJECT_ID":     fieldname.ProjectID,
	"DESCOPE_MANAGEMENT_KEY": fieldname.ManagementKey,
}
