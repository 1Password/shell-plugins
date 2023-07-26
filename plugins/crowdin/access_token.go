package crowdin

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.AccessToken,
		DocsURL: sdk.URL("https://developer.crowdin.com/configuration-file/#api-credentials-from-environment-variables"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Crowdin.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 80,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.ProjectID,
				MarkdownDescription: "Project ID used to authenticate to Crowdin.",
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Digits: true,
					},
				},
			},
			{
				Name:                fieldname.HostAddress,
				MarkdownDescription: "Base URL (for Crowdin Enterprise)",
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
					Prefix: "https://",
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CROWDIN_PERSONAL_TOKEN": fieldname.Token,
	"CROWDIN_PROJECT_ID":     fieldname.OrgID,
	"CROWDIN_BASE_URL":       fieldname.HostAddress,
}
