package tofu

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "tofu",
		Platform: schema.PlatformInfo{
			Name:     "OpenTofu",
			Homepage: sdk.URL("https://opentofu.org"),
		},
		Executables: []schema.Executable{
			TofuCLI(),
		},
	}
}
