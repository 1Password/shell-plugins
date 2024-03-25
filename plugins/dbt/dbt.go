package dbt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func dbtCLI() schema.Executable {
	return schema.Executable{
		Name:    "dbt",
		Runs:    []string{"dbt"},
		DocsURL: sdk.URL("https://docs.getdbt.com/docs/core/connect-data-platform/about-core-connections"),
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
