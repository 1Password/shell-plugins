package docker

import (
	"context"

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
		DocsURL:       sdk.URL("https://docs.docker.com/docker-hub/access-tokens/"),
		ManagementURL: sdk.URL("https://docs.docker.com/docker-hub/access-tokens/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Docker.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 18,
					Charset: schema.Charset{
						Lowercase: true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryDockerConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DOCKERHUB_TOKEN": fieldname.Token, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Personal Access Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryDockerConfigFile() sdk.Importer {
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
