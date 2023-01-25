package nprev

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NixPkgsCLI() schema.Executable {
	return schema.Executable{
		Name:    "nixpkgs-review CLI",
		Runs:    []string{"nixpkgs-review"},
		DocsURL: sdk.URL("https://github.com/Mic92/nixpkgs-review"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
