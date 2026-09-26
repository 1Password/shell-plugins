package apono

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AponoCLI() schema.Executable {
	return schema.Executable{
		Name:    "Apono CLI",
		Runs:    []string{"apono"},
		DocsURL: sdk.URL("https://docs.apono.io/docs/access-requests-and-approvals/cli/install-and-manage-the-apono-cli"),
		NeedsAuth: needsauth.IfAll(
			// The Apono CLI only accepts a token on "apono login"; all other
			// commands authenticate through the session it stores itself.
			needsauth.ForCommand("login"),
			needsauth.NotForHelp(),
			needsauth.NotWhenContainsArgs("--personal-token"),
			needsauth.NotWhenContainsArgs("__complete"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAPIToken,
			},
		},
	}
}
