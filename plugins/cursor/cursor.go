package cursor

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CursorCLI() schema.Executable {
	return schema.Executable{
		Name:    "Cursor CLI",
		Runs:    []string{"agent"},
		DocsURL: sdk.URL("https://cursor.com/docs/cli/overview"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("login"),
			needsauth.NotForExactArgs("logout"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
