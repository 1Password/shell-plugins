package sops

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sops",
		Platform: schema.PlatformInfo{
			Name:     "SOPS",
			Homepage: sdk.URL("https://github.com/getsops/sops"),
		},
		Credentials: []schema.CredentialType{
			AgeSecretKey(),
		},
		Executables: []schema.Executable{
			SOPSCLI(),
			HelmCLI(),
		},
	}
}
