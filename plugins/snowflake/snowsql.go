package snowflake

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SnowflakeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Snowflake",
		Runs:    []string{"snowsql"},
		DocsURL: sdk.URL("https://docs.snowflake.com/en/user-guide/snowsql.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("-a"),
			needsauth.NotWhenContainsArgs("--accountname"),
			needsauth.NotWhenContainsArgs("-u"),
			needsauth.NotWhenContainsArgs("--username"),
			needsauth.NotWhenContainsArgs("--authenticator"),
			needsauth.NotWhenContainsArgs("--config"),
			needsauth.NotWhenContainsArgs("-P"),
			needsauth.NotWhenContainsArgs("--prompt"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
