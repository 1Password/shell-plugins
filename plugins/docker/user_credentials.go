package docker

import (
	"context"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          "",                                                         // TODO: Register name in project://sdk/schema/credname/names.go
		DocsURL:       sdk.URL("https://docker.com/docs/user_credentials"),        // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.docker.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: " used to authenticate to Docker.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: " used to authenticate to Docker.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.TempFile(dockerConfig, provision.AtFixedPath("~/.docker/config.json")),
		Importer: importer.TryAll(
			TryDockerConfigFile(),
		)}
}

type PerHostAuth struct {
	UsernamePassword string `json:"auth"`
}

type Config struct {
	Auths map[string]PerHostAuth `json:"auths"`
}

// TODO: Check if the platform stores the User Credentials in a local config file, and if so,
// implement the function below to add support for importing it.
func TryDockerConfigFile() sdk.Importer {
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
