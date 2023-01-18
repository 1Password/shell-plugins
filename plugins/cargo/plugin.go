package cargo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cargo",
		Platform: schema.PlatformInfo{
			Name:     "Cargo",
			Homepage: sdk.URL("https://crates.io"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			CargoCLI(),
		},
	}
}
