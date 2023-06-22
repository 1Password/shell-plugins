package appwrite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AppwriteCLI() schema.Executable {
	return schema.Executable{
		Name:    "Appwrite CLI",
		Runs:    []string{"appwrite"},
		DocsURL: sdk.URL("https://appwrite.io/docs/command-line"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.CLIToken,
			},
		},
	}
}
