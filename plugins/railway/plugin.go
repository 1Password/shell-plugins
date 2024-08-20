package railway

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "railway",
		Platform: schema.PlatformInfo{
			Name:     "Railway",
			Homepage: sdk.URL("https://railway.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			RailwayCLI(),
		},
	}
}
