package fooctl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "fooctl",
		Platform: schema.PlatformInfo{
			Name:     "Fooctl",
			Homepage: sdk.URL("https://localhost/fooctl-cli"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			FooctlCLI(),
		},
	}
}
