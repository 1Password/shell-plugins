package laravelvapor

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LaravelVaporCLI() schema.Executable {
	return schema.Executable{
		Name:    "Laravel Vapor CLI",
		Runs:    []string{"vapor"},
		DocsURL: sdk.URL("https://docs.vapor.build/1.0/introduction.html#installing-the-vapor-cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("login"),
			needsauth.NotForExactArgs("list"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
