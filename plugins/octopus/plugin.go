package octopus

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "octopus",
		Platform: schema.PlatformInfo{
			Name:     "Octopus Deploy",
			Homepage: sdk.URL("https://octopus.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			OctopusDeployCLI(),
		},
	}
}
