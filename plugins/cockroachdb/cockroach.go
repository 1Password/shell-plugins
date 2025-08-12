package cockroachdb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Cockroach() schema.Executable {
	return schema.Executable{
		Name:    "cockroach",
		Runs:    []string{"cockroach"},
		DocsURL: sdk.URL("https://www.cockroachlabs.com/docs/stable/cockroach-sql.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.ForCommand("sql"),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
