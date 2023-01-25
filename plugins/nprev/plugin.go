package nprev

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "nprev",
		Platform: schema.PlatformInfo{
			Name:     "NixPkgs",
			Homepage: sdk.URL("https://github.com/Mic92/nixpkgs-review"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			NixPkgsCLI(),
		},
	}
}
