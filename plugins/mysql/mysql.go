package mysql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func Mysql() schema.Executable {
	return schema.Executable{
		Name:      "mysql",
		Runs:      []string{"mysql"},
		DocsURL:   sdk.URL("https://dev.mysql.com/doc/refman/en/mysql.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.DatabaseCredentials,
				Provisioner: provision.TempFile(mysqlConfig, provision.Filename("my.cnf"), provision.AddArgs("--defaults-file={{ .Path }}")),
			},
		},
	}
}
