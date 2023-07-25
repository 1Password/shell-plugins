package postgresql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "postgresql",
		Platform: schema.PlatformInfo{
			Name:     "PostgreSQL",
			Homepage: sdk.URL("https://postgresql.org"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			Psql(),
			Pg_dump(),
			Pg_restore(),
		},
	}
}
