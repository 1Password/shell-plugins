package postman

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
		DocsURL:       sdk.URL("https://learning.postman.com/docs/developer/postman-api/intro-api/#generating-a-postman-api-key"),
		ManagementURL: sdk.URL("https://web.postman.co/settings/me/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to postman.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 64,
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
			TryPostmanConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"POSTMAN_API_KEY": fieldname.APIKey,
}

func TryPostmanConfigFile() sdk.Importer {
	return importer.TryFile("~/.postman/postmanrc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		var defaultProfile *PostmanProfile

		for _, profile := range config.Login.Profiles {
			if profile.Alias == "default" {
				defaultProfile = &profile
				break
			}
		}

		if defaultProfile == nil || defaultProfile.PostmanAPIKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: defaultProfile.PostmanAPIKey,
			},
		})
	})
}

type PostmanProfile struct {
	Alias         string `json:"alias"`
	PostmanAPIKey string `json:"postmanApiKey"`
	Username      string `json:"username"`
}

type Config struct {
	Login struct {
		Profiles []PostmanProfile `json:"_profiles"`
	} `json:"login"`
	Updates struct {
		UpdateCheckedTimestamp int64 `json:"updateCheckedTimestamp"`
	} `json:"updates"`
}
