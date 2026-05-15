package postgresql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Pgcli() schema.Executable {
	return schema.Executable{
		Name:      "pgcli",
		Runs:      []string{"pgcli"},
		DocsURL:   sdk.URL("https://pgcli.com/docs"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
