package fastly

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "fastly",
		Platform: schema.PlatformInfo{
			Name:     "Fastly",
			Homepage: sdk.URL("https://fastly.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			FastlyCLI(),
		},
	}
}
