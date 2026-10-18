package netlify

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
		DocsURL:       sdk.URL("https://docs.netlify.com/api-and-cli-guides/cli-guides/get-started-with-cli"),
		ManagementURL: sdk.URL("https://app.netlify.com/user/applications#personal-access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Netlify.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.MacOnly(TryNetlifyConfigFile("~/Library/Preferences/netlify/config.json")),
			importer.LinuxOnly(TryNetlifyConfigFile("~/.config/netlify/config.json")),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NETLIFY_AUTH_TOKEN": fieldname.Token,
}

// TryNetlifyConfigFile imports the token that `netlify login` stores in the global config.json,
// under users.<userId>.auth.token for the active userId.
func TryNetlifyConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		user, ok := config.Users[config.UserID]
		if !ok || user.Auth.Token == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: user.Auth.Token,
			},
		})
	})
}

type Config struct {
	UserID string                `json:"userId"`
	Users  map[string]ConfigUser `json:"users"`
}

type ConfigUser struct {
	Auth struct {
		Token string `json:"token"`
	} `json:"auth"`
}
