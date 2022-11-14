package mysql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Mysql() schema.Executable {
	return schema.Executable{
		Runs:      []string{"mysql"},
		Name:      "mysql",
		DocsURL:   sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/mysql.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
	}
}
