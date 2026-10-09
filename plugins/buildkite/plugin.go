package buildkite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "buildkite",
		Platform: schema.PlatformInfo{
			Name:     "Buildkite",
			Homepage: sdk.URL("https://buildkite.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			BuildkiteCLI(),
		},
	}
}
