package junie

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func JunieCLI() schema.Executable {
	return schema.Executable{
		Name:    "JetBrains Junie CLI",
		Runs:    []string{"junie"},
		DocsURL: sdk.URL("https://junie.jetbrains.com/docs/junie-cli-usage.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
