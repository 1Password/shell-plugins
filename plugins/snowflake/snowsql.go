package snowflake

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SnowflakeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Snowflake",
		Runs:      []string{"snowsql"},
		DocsURL:   sdk.URL("https://docs.snowflake.com/en/user-guide/snowsql.html"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
