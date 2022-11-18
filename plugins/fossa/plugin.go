package fossa

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "fossa",
		Platform: schema.PlatformInfo{
			Name:     "FOSSA",
			Homepage: sdk.URL("https://fossa.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			FOSSACLI(),
		},
	}
}
