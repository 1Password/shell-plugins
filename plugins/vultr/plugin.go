package vultr

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "vultr",
		Platform: schema.PlatformInfo{
			Name:     "Vultr",
			Homepage: sdk.URL("https://vultr.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			VultrCLI(),
		},
	}
}
