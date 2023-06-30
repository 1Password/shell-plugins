package cratedb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cratedb",
		Platform: schema.PlatformInfo{
			Name:     "CrateDB",
			Homepage: sdk.URL("https://crate.io/"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			CrateDBCLI(),
		},
	}
}
