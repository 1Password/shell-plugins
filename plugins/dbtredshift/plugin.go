package dbtredshift

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "dbtredshift",
		Platform: schema.PlatformInfo{
			Name:     "DBT Redshift",
			Homepage: sdk.URL("https://www.getdbt.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			dbtredshiftCLI(),
		},
	}
}
