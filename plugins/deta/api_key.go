package deta

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
		DocsURL:       sdk.URL("https://deta.space/docs/en"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://deta.space"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Deta.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 25,
					Prefix: "bE53FsM6_", // TODO: Check if this is correct
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
			TryDetaConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	// "DETA_API_KEY": fieldname.APIKey, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the API Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryDetaConfigFile() sdk.Importer {
	return importer.TryFile("~/.detaspace/space_tokens", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
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
