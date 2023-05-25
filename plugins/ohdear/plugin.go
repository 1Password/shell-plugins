package ohdear

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ohdear",
		Platform: schema.PlatformInfo{
			Name:     "Oh Dear",
			Homepage: sdk.URL("https://ohdear.app"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			OhDearCLI(),
		},
	}
}
