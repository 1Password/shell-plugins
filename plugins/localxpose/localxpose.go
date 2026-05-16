package localxpose

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LocalXposeCLI() schema.Executable {
	return schema.Executable{
		Name:    "LocalXpose CLI",
		Runs:    []string{"loclx"},
		DocsURL: sdk.URL("https://localxpose.io/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("account", "login"),     // skip 1Password authentication for "loclx account login" and its subcommands
			needsauth.NotWhenContainsArgs("a", "login"),           // skip 1Password authentication for "loclx account login" and its subcommands
			needsauth.NotWhenContainsArgs("account", "logout"),    // skip 1Password authentication for "loclx account logout" and its subcommands
			needsauth.NotWhenContainsArgs("a", "logout"),          // skip 1Password authentication for "loclx account logout" and its subcommands
			needsauth.NotWhenContainsArgs("service", "restart"),   // skip 1Password authentication for "loclx service restart" and its subcommands
			needsauth.NotWhenContainsArgs("service", "start"),     // skip 1Password authentication for "loclx service start" and its subcommands
			needsauth.NotWhenContainsArgs("service", "status"),    // skip 1Password authentication for "loclx service status" and its subcommands
			needsauth.NotWhenContainsArgs("service", "stop"),      // skip 1Password authentication for "loclx service stop" and its subcommands
			needsauth.NotWhenContainsArgs("service", "uninstall"), // skip 1Password authentication for "loclx service uninstall" and its subcommands
			needsauth.NotWhenContainsArgs("setting"),              // skip 1Password authentication for "loclx setting" and its subcommands
			needsauth.NotWhenContainsArgs("update"),               // skip 1Password authentication for "loclx update" and its subcommands
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
