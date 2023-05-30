package flyctl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "flyctl",
		Platform: schema.PlatformInfo{
			Name:     "Fly.io",
			Homepage: sdk.URL("https://fly.io"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			FlyctlCLI(),
			FlyCLI(),
		},
	}
}
