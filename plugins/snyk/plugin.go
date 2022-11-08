package snyk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "snyk",
		Platform: schema.PlatformInfo{
			Name:     "Snyk",
			Homepage: sdk.URL("https://snyk.io"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			SnykCLI(),
		},
	}
}
