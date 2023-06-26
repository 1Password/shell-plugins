package caprover

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "caprover",
		Platform: schema.PlatformInfo{
			Name:     "Caprover",
			Homepage: sdk.URL("https://caprover.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			CaproverCLI(),
		},
	}
}
