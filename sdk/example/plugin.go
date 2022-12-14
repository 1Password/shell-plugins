package example

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "example",
		Platform: schema.PlatformInfo{
			Name:     "Example",
			Homepage: sdk.URL("https://example.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			ExampleCLI(),
		},
	}
}
