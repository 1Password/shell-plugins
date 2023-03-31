package akamai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "akamai",
		Platform: schema.PlatformInfo{
			Name:     "Akamai",
			Homepage: sdk.URL("https://akamai.com"),
		},
		Credentials: []schema.CredentialType{
			APIClientCredentials(),
		},
		Executables: []schema.Executable{
			AkamaiCLI(),
		},
	}
}
