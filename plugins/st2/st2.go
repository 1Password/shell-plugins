package st2

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func StackStormCLI() schema.Executable {
	return schema.Executable{
		Name:      "StackStorm CLI",
		Runs:      []string{"st2"},
		DocsURL:   sdk.URL("https://docs.stackstorm.com/reference/cli.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.UserPass,
			},
		},
	}
}
