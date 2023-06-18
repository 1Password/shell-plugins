package nirmata

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "nirmata",
		Platform: schema.PlatformInfo{
			Name:     "Nirmata",
			Homepage: sdk.URL("https://nirmata.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			NirmataCLI(),
		},
	}
}
