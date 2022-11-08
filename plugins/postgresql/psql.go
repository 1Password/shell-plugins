package postgresql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Psql() schema.Executable {
	return schema.Executable{
		Runs:      []string{"psql"},
		Name:      "psql",
		DocsURL:   sdk.URL("https://www.postgresql.org/docs/current/app-psql.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
	}
}
