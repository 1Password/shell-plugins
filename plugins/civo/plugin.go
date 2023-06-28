package civo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "civo",
		Platform: schema.PlatformInfo{
			Name:     "Civo",
			Homepage: sdk.URL("https://civo.com"), 
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			CivoCLI(),
		},
	}
}
