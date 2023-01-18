package cachix

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cachix",
		Platform: schema.PlatformInfo{
			Name:     "Cachix",
			Homepage: sdk.URL("https://www.cachix.org"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			CachixCLI(),
		},
	}
}
