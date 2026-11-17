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
			needsauth.NotForCommand("account", "login"),     // skip 1Password authentication for "loclx account login" and its subcommands
			needsauth.NotForCommand("a", "login"),           // skip 1Password authentication for "loclx account login" and its subcommands
			needsauth.NotForCommand("account", "logout"),    // skip 1Password authentication for "loclx account logout" and its subcommands
			needsauth.NotForCommand("a", "logout"),          // skip 1Password authentication for "loclx account logout" and its subcommands
			needsauth.NotForCommand("service", "restart"),   // skip 1Password authentication for "loclx service restart" and its subcommands
			needsauth.NotForCommand("service", "start"),     // skip 1Password authentication for "loclx service start" and its subcommands
			needsauth.NotForCommand("service", "status"),    // skip 1Password authentication for "loclx service status" and its subcommands
			needsauth.NotForCommand("service", "stop"),      // skip 1Password authentication for "loclx service stop" and its subcommands
			needsauth.NotForCommand("service", "uninstall"), // skip 1Password authentication for "loclx service uninstall" and its subcommands
			needsauth.NotForCommand("setting"),              // skip 1Password authentication for "loclx setting" and its subcommands
			needsauth.NotForCommand("update"),               // skip 1Password authentication for "loclx update" and its subcommands
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
