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
		DocsURL:       sdk.URL("https://docs.netlify.com/cli/get-started/#authentication"),
		ManagementURL: sdk.URL("https://app.netlify.com/user/applications#personal-access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Netlify.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 43,
					Prefix: "tGtp-",
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
			TryNetlifyConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName {
	"NETLIFY_AUTH_TOKEN": fieldname.Token,
}

func TryNetlifyConfigFile() sdk.Importer {
	return importer.TryFile("~/Library/Preferences/netlify/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.Users != nil {
			for _, user := range config.Users {
				if user.Auth != nil && user.Auth.Token != "" {
					out.AddCandidate(sdk.ImportCandidate{
						Fields: map[sdk.FieldName]string{
							fieldname.Token: user.Auth.Token,
						},
					})
				}
			}
		}
	})
}

type Config struct {
	TelemetryDisabled bool                `json:"telemetryDisabled"`
	CliID             string              `json:"cliId"`
	UserID            string              `json:"userId"`
	Users             map[string]UserInfo `json:"users"`
}

// UserInfo represents the user information in the config file
type UserInfo struct {
	ID    string        `json:"id"`
	Name  string        `json:"name"`
	Email string        `json:"email"`
	Auth  *UserAuthInfo `json:"auth"`
}

// UserAuthInfo represents the authentication information of a user
type UserAuthInfo struct {
	Token   string `json:"token"`
	Github  struct{} `json:"github"` // Empty struct for placeholder, you can add additional fields if needed
}
