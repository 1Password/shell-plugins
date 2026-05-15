package postgresql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Pg_dump() schema.Executable {
	return schema.Executable{
		Name:      "pg_dump",
		Runs:      []string{"pg_dump"},
		DocsURL:   sdk.URL("https://www.postgresql.org/docs/current/app-pgdump.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
