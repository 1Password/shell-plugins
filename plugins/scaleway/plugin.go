package scaleway

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "scaleway",
		Platform: schema.PlatformInfo{
			Name:     "Scaleway",
			Homepage: sdk.URL("https://scaleway.com"),
		},
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
		Executables: []schema.Executable{
			ScalewayCLI(),
		},
	}
}
