package appwrite

import (
	"context"

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
		DocsURL:       sdk.URL("https://appwrite.com/docs/api_key"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.appwrite.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Appwrite.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 20,
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
			TryAppwriteConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"APPWRITE_API_KEY": fieldname.APIKey, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the API Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryAppwriteConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.APIKey == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.APIKey: config.APIKey,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	APIKey string
// }
