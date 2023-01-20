package snowflake

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SnowflakeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Snowflake CLI", // TODO: Check if this is correct
		Runs:      []string{"snowsql"},
		DocsURL:   sdk.URL("https://snowflake.com/docs/cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
