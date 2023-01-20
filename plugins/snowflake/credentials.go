package snowflake

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credentials,
		DocsURL:       sdk.URL("https://snowflake.com/docs/credentials"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.snowflake.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.,
				MarkdownDescription: " used to authenticate to Snowflake.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 11,
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
			TrySnowflakeConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SNOWFLAKE": fieldname., // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Credentials in a local config file, and if so,
// implement the function below to add support for importing it.
func TrySnowflakeConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config. == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.: config.,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	 string
// }
