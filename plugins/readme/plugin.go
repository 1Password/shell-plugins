package readme

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "readme",
		Platform: schema.PlatformInfo{
			Name:     "ReadMe",
			Homepage: sdk.URL("https://readme.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ReadMeCLI(),
		},
	}
}
