package heroku

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HerokuCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"heroku"},
		Name:      "Heroku CLI",
		DocsURL:   sdk.URL("https://devcenter.heroku.com/articles/heroku-cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
