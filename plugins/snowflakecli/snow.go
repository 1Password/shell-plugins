package snowflakecli

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SnowflakeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Snowflake CLI",
		Runs:      []string{"snow"},
		DocsURL:   sdk.URL("https://docs.snowflake.com/en/developer-guide/snowflake-cli/index"),
		NeedsAuth: snowNeedsAuth(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}

func snowNeedsAuth() sdk.NeedsAuthentication {
	return needsauth.IfAll(
		needsauth.NotForHelpOrVersion(),
		needsauth.NotWithoutArgs(),
		// The user explicitly chose which connection or config file to use.
		needsauth.NotWhenContainsArgs("-c"),
		needsauth.NotWhenContainsArgs("--connection"),
		needsauth.NotWhenContainsArgs("-x"),
		needsauth.NotWhenContainsArgs("--temporary-connection"),
		needsauth.NotWhenContainsArgs("--config-file"),
		// Commands that don't connect to Snowflake. These don't accept the -x flag added by the provisioner.
		needsauth.NotForCommand("app", "bundle"),
		needsauth.NotForCommand("auth", "list-connections"),
		needsauth.NotForCommand("auth", "oidc", "read-token"),
		needsauth.NotForCommand("connection", "add"),
		needsauth.NotForCommand("connection", "list"),
		needsauth.NotForCommand("connection", "remove"),
		needsauth.NotForCommand("connection", "set-default"),
		needsauth.NotForCommand("custom-image", "validate"),
		needsauth.NotForCommand("helpers"),
		needsauth.NotForCommand("init"),
		needsauth.NotForCommand("plugin"),
		needsauth.NotForCommand("upgrade"),
	)
}
