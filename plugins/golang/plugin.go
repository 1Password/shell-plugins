package golang

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "golang",
		Platform: schema.PlatformInfo{
			Name:     "Go",
			Homepage: sdk.URL("https://go.dev/"),
		},
		Executables: []schema.Executable{
			GoCLI(),
		},
	}
}
