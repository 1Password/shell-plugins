package stackhero

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func StackheroCLI() schema.Executable {
	return schema.Executable{
		Name:    "Stackhero CLI",
		Runs:    []string{"stackhero"},
		DocsURL: sdk.URL("https://www.stackhero.io/stackhero/documentations/Use-the-CLI"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForCommand("login"),
			needsauth.NotForCommand("logout"),
			needsauth.NotForCommand("self-update"),
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("--help-agents"),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
