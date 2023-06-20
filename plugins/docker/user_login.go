package docker

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserLogin() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.UserLogin,
		DocsURL:       sdk.URL("https://docker.com/docs/user_login"),              // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.docker.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username of the user registered in Docker.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password of the user registered in Docker.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(

			TryDockerConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DOCKER_USERNAME": fieldname.Username,
	"DOCKER_PASSWORD": fieldname.Password,
}

func TryDockerConfigFile() sdk.Importer {
	return importer.TryFile("~/.docker/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {

		var config struct {
			Auths map[string]struct {
				Auth string `json:"auth"`
			} `json:"auths"`
		}

		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}
		for _, auth := range config.Auths {
			decoded, err := base64.StdEncoding.DecodeString(auth.Auth)
			if err != nil {
				out.AddError(err)
				return
			}
			credentials := string(decoded)
			username, password := parseCredentials(credentials)
			if username != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: username,
						fieldname.Password: password,
					},
				})
			}
		}

	})
}

func parseCredentials(credentials string) (username string, password string) {
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) == 2 {
		username = parts[0]
		password = parts[1]
	}
	return username, password
}
