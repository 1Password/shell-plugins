package tunnelto

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "tunnelto",
		Platform: schema.PlatformInfo{
			Name:     "tunnelto.dev",
			Homepage: sdk.URL("https://tunnelto.dev"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			tunneltodevCLI(),
		},
	}
}
