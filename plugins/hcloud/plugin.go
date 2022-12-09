package hcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "hcloud",
		Platform: schema.PlatformInfo{
			Name:     "Hetzner Cloud",
			Homepage: sdk.URL("https://console.hetzner.cloud"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			HetznerCloudCLI(),
		},
	}
}
