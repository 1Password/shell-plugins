package atlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func MongoshCLI() schema.Executable {
	return schema.Executable{
		Name:    "MongoDB Shell",
		Runs:    []string{"mongosh"},
		DocsURL: sdk.URL("https://www.mongodb.com/docs/mongodb-shell/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.DatabaseCredentials,
				Provisioner: provision.ArgsAtIndex(indexToProvisionAt, argsToProvision),
				Importer:    nil,
			},
		},
	}
}
