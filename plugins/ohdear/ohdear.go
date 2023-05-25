package ohdear

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OhDearCLI() schema.Executable {
	return schema.Executable{
		Name:    "Oh Dear CLI",
		Runs:    []string{"ohdear"},
		DocsURL: sdk.URL("https://ohdear.app/docs/integrations/our-cli-tool"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
