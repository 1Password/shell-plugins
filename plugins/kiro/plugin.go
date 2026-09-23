package kiro

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "kiro",
		Platform: schema.PlatformInfo{
			Name:     "Kiro",
			Homepage: sdk.URL("https://kiro.dev"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			KiroCLI(),
		},
	}
}
