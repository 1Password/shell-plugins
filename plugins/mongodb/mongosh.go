package mongodb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func MongoshCLI() schema.Executable {
	return schema.Executable{
		Name:    "MongoDB Shell",
		Runs:    []string{"mongosh"},
		DocsURL: sdk.URL("https://www.mongodb.com/docs/mongodb-shell/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotWhenContainsArgs("--host"),     // skip 1Password auth if host is provided
			needsauth.NotWhenContainsArgs("--port"),     // skip 1Password auth if port is provided
			needsauth.NotWhenContainsArgs("--username"), // skip 1Password auth if username is provided
			needsauth.NotWhenContainsArgs("--password"), // skip 1Password auth if password is provided
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.DatabaseCredentials,
				Provisioner: mongodbShellProvisioner(),
			},
		},
	}
}
