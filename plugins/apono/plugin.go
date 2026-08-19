package apono

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "apono",
		Platform: schema.PlatformInfo{
			Name:     "Apono",
			Homepage: sdk.URL("https://apono.io"),
		},
		Credentials: []schema.CredentialType{
			PersonalAPIToken(),
		},
		Executables: []schema.Executable{
			AponoCLI(),
		},
	}
}
