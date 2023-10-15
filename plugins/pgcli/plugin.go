package pgcli

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "pgcli",
		Platform: schema.PlatformInfo{
			Name:     "PostgreSQL",
			Homepage: sdk.URL("https://pgcli.com"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			PostgreSQLCLI(),
		},
	}
}
