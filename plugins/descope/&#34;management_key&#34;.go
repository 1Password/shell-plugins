package descope

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func &#34;ManagementKey&#34;() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.&#34;ManagementKey&#34;, // TODO: Register name in project://sdk/schema/credname/names.go
		DocsURL:       sdk.URL("https://descope.com/docs/&#34;management_key&#34;"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.descope.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key&#34;,
				MarkdownDescription: "Key&#34; used to authenticate to Desope.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 71,
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
			TryDesopeConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DESCOPE_KEY&#34;": fieldname.Key&#34;, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the &#34;Management Key&#34; in a local config file, and if so,
// implement the function below to add support for importing it.
func TryDesopeConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Key&#34; == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Key&#34;: config.Key&#34;,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Key&#34; string
// }
