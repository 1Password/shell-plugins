package aiven

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "aiven",
		Platform: schema.PlatformInfo{
			Name:     "Aiven.io",
			Homepage: sdk.URL("https://aiven.io"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			AivenCLI(),
		},
	}
}
