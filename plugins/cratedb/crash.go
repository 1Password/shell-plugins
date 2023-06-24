package cratedb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CrateDBCLI() schema.Executable {
	return schema.Executable{
		Name:      "CrateDB Shell", // TODO: Check if this is correct
		Runs:      []string{"crash"},
		DocsURL:   sdk.URL("https://crate.io/docs/crate/crash/en/latest/"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
