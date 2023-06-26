package deta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "deta",
		Platform: schema.PlatformInfo{
			Name:     "Deta",
			Homepage: sdk.URL("https://deta.space"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			DetaCLI(),
		},
	}
}
