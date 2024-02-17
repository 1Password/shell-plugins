package wiz

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "wiz",
		Platform: schema.PlatformInfo{
			Name:     "Wiz",
			Homepage: sdk.URL("https://wiz.io"),
		},
		Credentials: []schema.CredentialType{
			SecretKey(),
		},
		Executables: []schema.Executable{
			WizCLI(),
		},
	}
}
