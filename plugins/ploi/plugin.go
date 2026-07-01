package ploi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ploi",
		Platform: schema.PlatformInfo{
			Name:     "Ploi CLI",
			Homepage: sdk.URL("https://ploi.io"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			PloiCLICLI(),
		},
	}
}
