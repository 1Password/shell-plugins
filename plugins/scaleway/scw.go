package scaleway

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ScalewayCLI() schema.Executable {
	return schema.Executable{
		Name:    "Scaleway CLI",
		Runs:    []string{"scw"},
		DocsURL: sdk.URL("https://www.scaleway.com/en/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("-c"),
			needsauth.NotWhenContainsArgs("--config"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
