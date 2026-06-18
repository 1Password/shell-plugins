package opencode

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "opencode",
		Platform: schema.PlatformInfo{
			Name:     "opencode",
			Homepage: sdk.URL("https://opencode.ai"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			OpencodeCLI(),
		},
	}
}
