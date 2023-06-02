package clilol

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "clilol",
		Platform: schema.PlatformInfo{
			Name:     "omg.lol",
			Homepage: sdk.URL("https://omg.lol"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			omglolCLI(),
		},
	}
}
