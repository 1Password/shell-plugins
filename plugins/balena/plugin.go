package balena

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "balena",
		Platform: schema.PlatformInfo{
			Name:     "Balena",
			Homepage: sdk.URL("https://www.balena.io"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			BalenaCLI(),
		},
	}
}
