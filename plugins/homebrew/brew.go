package homebrew

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HomebrewCLI() schema.Executable {
	return schema.Executable{
		Name:    "Homebrew CLI",
		Runs:    []string{"brew"},
		DocsURL: sdk.URL("https://brew.sh/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.IfAny(
				needsauth.ForCommand("audit"),
				needsauth.ForCommand("bump"),
				needsauth.ForCommand("bump-cask-pr"),
				needsauth.ForCommand("bump-formula-pr"),
				needsauth.ForCommand("install"),
				needsauth.ForCommand("reinstall"),
				needsauth.ForCommand("search"),
				needsauth.ForCommand("update"),
				needsauth.ForCommand("upgrade"),
			),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
