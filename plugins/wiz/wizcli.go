package wiz

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func WizCLI() schema.Executable {
	return schema.Executable{
		Name:    "Wiz CLI",
		Runs:    []string{"wizcli"},
		DocsURL: sdk.URL("https://docs.wiz.io/wiz-docs/docs/wiz-cli-overview"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
