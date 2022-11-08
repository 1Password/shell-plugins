package stripe

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "stripe",
		Platform: schema.PlatformInfo{
			Name:     "Stripe",
			Homepage: sdk.URL("https://stripe.com"),
		},
		Credentials: []schema.CredentialType{
			SecretKey(),
		},
		Executables: []schema.Executable{
			StripeCLI(),
		},
	}
}
