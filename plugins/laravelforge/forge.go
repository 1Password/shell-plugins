package laravelforge

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LaravelForgeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Laravel Forge CLI",
		Runs:    []string{"forge"},
		DocsURL: sdk.URL("https://forge.laravel.com/docs/1.0/cli.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("login"),
			needsauth.NotForExactArgs("logout"),
			needsauth.NotForExactArgs("list"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
