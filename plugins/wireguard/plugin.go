package wireguard

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "wireguard",
		Platform: schema.PlatformInfo{
			Name:     "Wireguard VPN",
			Homepage: sdk.URL("https://wireguard.com"),
		},
		Credentials: []schema.CredentialType{
			AccessConfig(),
		},
		Executables: []schema.Executable{
			WireguardVPNCLI(),
		},
	}
}
