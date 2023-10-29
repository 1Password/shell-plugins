package shopify

import (
	"context"

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
		DocsURL:       sdk.URL("https://shopify.com/docs/cli_token"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.shopify.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Shopify.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 69,
					Prefix: "atkn_", // TODO: Check if this is correct
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
			TryShopifyConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SHOPIFY_TOKEN": fieldname.Token, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the CLI Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryShopifyConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Token == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Token: config.Token,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Token string
// }
