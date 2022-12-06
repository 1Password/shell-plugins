package github

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"),
		ManagementURL: sdk.URL("https://github.com/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to GitHub.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 93,
					Prefix: "github_pat_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The GitHub host to authenticate to. Defaults to 'github.com'.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryAllEnvVars(fieldname.Token, "GH_TOKEN", "GITHUB_PAT"),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":             fieldname.Host,
				"GH_ENTERPRISE_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":                 fieldname.Host,
				"GITHUB_ENTERPRISE_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":  fieldname.Host,
				"GH_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":      fieldname.Host,
				"GITHUB_TOKEN": fieldname.Token,
			}),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"GITHUB_TOKEN": fieldname.Token,
}
