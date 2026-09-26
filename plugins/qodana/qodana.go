package qodana

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func QodanaCLI() schema.Executable {
	return schema.Executable{
		Name:    "Qodana CLI",
		Runs:    []string{"qodana"},
		DocsURL: sdk.URL("https://www.jetbrains.com/help/qodana/quick-start.html#quickstart-prerequisites-qodana-cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("cloc"),
			needsauth.NotWhenContainsArgs("completion"),
			needsauth.NotWhenContainsArgs("contributors"),
			needsauth.NotWhenContainsArgs("pull"),
			needsauth.NotWhenContainsArgs("show"),
			needsauth.NotWhenContainsArgs("view"),
		),
		Uses: []schema.CredentialUsage{
			{Name: credname.APIToken},
		},
	}
}
