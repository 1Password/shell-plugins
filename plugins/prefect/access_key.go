package prefect

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
		DocsURL:       sdk.URL("https://docs.prefect.io/2.10.13/cloud/users/api-keys/"),
		ManagementURL: sdk.URL("https://app.prefect.cloud/my/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key,
				MarkdownDescription: "Key used to authenticate to Prefect.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "pnu_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "URL for the Prefect workspace.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 124,
					Prefix: "https://api.prefect.cloud/api/accounts/",
					Charset: schema.Charset{
						Uppercase: false,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryPrefectConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"PREFECT_API_KEY": fieldname.Key,
	"PREFECT_API_URL": fieldname.URL,
}

// TODO: Check if the platform stores the Access Key in a local config file, and if so,
// implement the function below to add support for importing it.
func TryPrefectConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Key == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Key: config.Key,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Key string
// }
