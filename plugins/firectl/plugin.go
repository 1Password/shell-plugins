package firectl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "firectl",
		Platform: schema.PlatformInfo{
			Name:     "Fireworks AI",
			Homepage: sdk.URL("https://fireworks.ai"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			FirectlCLI(),
		},
	}
}
