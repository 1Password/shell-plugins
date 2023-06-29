package firebase

import (
	"context"

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
		DocsURL:       sdk.URL("https://firebase.google.com/docs/cli#cli-ci-systems"),
		ManagementURL: sdk.URL("https://firebase.google.com/docs/cli#cli-ci-systems"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to firebase.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 53,
					Prefix: "dummy_firebase_", // TODO: Check if this is correct
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
			TryfirebaseConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"FIREBASE_TOKEN": fieldname.Token, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Access Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryfirebaseConfigFile() sdk.Importer {
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
