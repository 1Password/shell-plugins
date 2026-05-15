package postgresql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Pg_restore() schema.Executable {
	return schema.Executable{
		Name:      "pg_restore",
		Runs:      []string{"pg_restore"},
		DocsURL:   sdk.URL("https://www.postgresql.org/docs/current/app-pgrestore.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
