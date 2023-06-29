package dbtredshift

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func dbtredshiftCLI() schema.Executable {
	return schema.Executable{
		Name:    "dbtredshift",
		Runs:    []string{"dbt"},
		DocsURL: sdk.URL("https://docs.getdbt.com/docs/core/connect-data-platform/redshift-setup"),
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
