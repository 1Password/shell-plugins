package shopify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ShopifyThemeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Shopify CLI",
		Runs:    []string{"shopify"},
		DocsURL: sdk.URL("https://github.com/Shopify/cli/blob/main/packages/cli/README.md#commands"),
		NeedsAuth: needsauth.IfAll(
			needsauth.ForCommand("theme"),
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("--password"),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.CLIToken,
			},
		},
	}
}
