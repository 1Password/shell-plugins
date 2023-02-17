package ngrok

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ngrokCLI() schema.Executable {
	return schema.Executable{
		Name:    "ngrok CLI",
		Runs:    []string{"ngrok"},
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/ngrok"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("--config"),   // skip 1Password authentication when any command contains "--config" flag
			needsauth.NotWhenContainsArgs("config"),     // skip 1Password authentication for "ngrok config" and its subcommands
			needsauth.NotWhenContainsArgs("update"),     // skip 1Password authentication for "ngrok update" and "ngrok update --channel=beta"
			needsauth.NotWhenContainsArgs("completion"), // skip 1Password authentication for "ngrok completion"
			needsauth.NotWhenContainsArgs("credits"),    // skip 1Password authentication for "ngrok credits"
			needsauth.NotWhenContainsArgs("service"),    // skip 1Password authentication for "ngrok service" subcommands because that setup involves defining/knowing in advance the config file path, which isn't something we can handle with 1Password Shell Plugins.
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
