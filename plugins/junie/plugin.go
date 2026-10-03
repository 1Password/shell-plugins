package junie

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "junie",
		Platform: schema.PlatformInfo{
			Name:     "JetBrains Junie",
			Homepage: sdk.URL("https://junie.jetbrains.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			JunieCLI(),
		},
	}
}
