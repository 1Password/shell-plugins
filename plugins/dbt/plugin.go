package dbt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "dbt",
		Platform: schema.PlatformInfo{
			Name:     "DBT",
			Homepage: sdk.URL("https://www.getdbt.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			dbtCLI(),
		},
	}
}
