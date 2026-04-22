package sourcegraph

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
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://sourcegraph.com/docs/cli"),
		ManagementURL: sdk.URL("https://sourcegraph.com/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Endpoint,
				AlternativeNames:    []string{"Website","URL"},
				MarkdownDescription: "Base URL for your Sourcegraph instance. Should start with https://",
				Composition: &schema.ValueComposition{
					Prefix:  "https://",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Token,
				AlternativeNames:    []string{"AccessToken"},
				MarkdownDescription: "Access token used to authenticate to Sourcegraph. Should start with sgp_",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length:  60,
					Prefix:  "sgp_",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'_'},
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SRC_ENDPOINT":     fieldname.Endpoint,
	"SRC_ACCESS_TOKEN": fieldname.Token,
}
