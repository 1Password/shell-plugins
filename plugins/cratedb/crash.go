package cratedb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CrateDBCLI() schema.Executable {
	return schema.Executable{
		Name:    "CrateDB Shell",
		Runs:    []string{"crash"},
		DocsURL: sdk.URL("https://crate.io/docs/crate/crash/en/latest/"),
		NeedsAuth: needsauth.needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.DatabaseCredentials,
				Provisioner: CrateArgsProvisioner{} ,
			},
		},
	}
}
