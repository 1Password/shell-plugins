package heroku

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
		DocsURL:       sdk.URL("https://devcenter.heroku.com/articles/authentication"),
		ManagementURL: sdk.URL("https://dashboard.heroku.com/account"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Heroku.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		Provisioner: provision.EnvVars(map[string]string{
			fieldname.APIKey: "HEROKU_API_KEY",
		}),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.APIKey, "HEROKU_API_KEY"),
			TryNetrcFile(),
		),
	}
}

// TryNetrcFile tries to find Heroku API keys in the ~/.netrc file
func TryNetrcFile() sdk.Importer {
	return importer.TryFile("~/.netrc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// TODO: Iterate over 'machine' entries to look for 'api.heroku.com' or 'git.heroku.com'
		// Example contents:
		//
		// machine api.heroku.com
		//   login me@example.com
		//   password c4cd94da15ea0544802c2cfd5ec4ead324327430
		// machine git.heroku.com
		//   login me@example.com
		//   password c4cd94da15ea0544802c2cfd5ec4ead324327430
	})
}
