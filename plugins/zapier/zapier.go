package zapier

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ZapierCLI() schema.Executable {
	return schema.Executable{
		Name:    "Zapier CLI",
		Runs:    []string{"zapier"},
		DocsURL: sdk.URL("https://platform.zapier.com/cli_docs/docs"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DeployKey,
			},
		},
	}
}
