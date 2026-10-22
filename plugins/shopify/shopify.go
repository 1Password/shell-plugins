package shopify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ShopifyCLI() schema.Executable {
	return schema.Executable{
		Name:    "Shopify CLI",
		Runs:    []string{"shopify"},
		DocsURL: sdk.URL("https://github.com/shopify/cli#readme"),
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
