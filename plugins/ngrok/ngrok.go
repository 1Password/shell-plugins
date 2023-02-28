package ngrok

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ngrokCLI() schema.Executable {

	var cmds []string
	cmds = append(cmds, "api")

	return schema.Executable{
		Name:    "ngrok CLI",
		Runs:    []string{"ngrok"},
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/ngrok"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("--config"),
			needsauth.NotForExactArgs("config"),
			needsauth.NotForExactArgs("update"),
		),

		Uses: []schema.CredentialUsage{
			{
				Name:    credname.AuthToken,
				Default: true,
			},
			{
				Name:     credname.APIKey,
				Commands: cmds,
			},
		},
	}
}
