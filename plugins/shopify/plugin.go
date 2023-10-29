package shopify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "shopify",
		Platform: schema.PlatformInfo{
			Name:     "Shopify",
			Homepage: sdk.URL("https://shopify.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			CLIToken(),
		},
		Executables: []schema.Executable{
			ShopifyCLI(),
		},
	}
}
