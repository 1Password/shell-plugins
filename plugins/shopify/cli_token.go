package shopify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func CLIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.CLIToken,
		DocsURL:       sdk.URL("https://shopify.dev/docs/apps/tools/cli/ci-cd#step-2-generate-a-cli-authentication-token"),
		ManagementURL: sdk.URL("https://partners.shopify.com/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Shopify.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 69,
					Prefix: "atkn_",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
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
	"SHOPIFY_TOKEN": fieldname.Token,
}
