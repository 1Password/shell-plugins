package nimble

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func nimCLI() schema.Executable {
	return schema.Executable{
		Name:    "nimble",
		Runs:    []string{"nimble"},
		DocsURL: sdk.URL("https://github.com/nim-lang/nimble"),
		NeedsAuth: needsauth.IfAll(
			needsauth.IfAny(
				needsauth.ForCommand("publish"),
			),
			needsauth.NotForHelp(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
