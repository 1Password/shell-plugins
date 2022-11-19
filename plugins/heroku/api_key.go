package heroku

import (
	"bufio"
	"context"
	"strings"

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
		s := bufio.NewScanner(strings.NewReader(string(contents)))
		var machine, login, password string
		for s.Scan() {
			if words := strings.Split(strings.Trim(s.Text(), " "), " "); len(words) >= 2 {
				switch words[0] {
				case "machine":
					if machine != "" {
						login, password = "", ""
					}
					machine = words[1]
				case "login":
					login = words[1]
				case "password":
					password = words[1]
				}
				if login != "" && password != "" && machine != "" {
					if machine == "api.heroku.com" || machine == "git.heroku.com" {
						out.AddCandidate(sdk.ImportCandidate{
							NameHint: login,
							Fields: map[string]string{
								fieldname.APIKey: password,
							},
						})
					}
					machine, login, password = "", "", ""
				}
			}
		}
	})
}
