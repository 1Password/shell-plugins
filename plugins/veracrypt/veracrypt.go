package veracrypt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"

)

func VeraCryptCLI() schema.Executable {
	return schema.Executable{
		Name:    "VeraCrypt CLI",
		Runs:    []string{"veracrypt"},
		DocsURL: sdk.URL("https://www.veracrypt.fr/en/Documentation.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: sdk.CredentialName("Volume Password"),
			},
		},
	}
}