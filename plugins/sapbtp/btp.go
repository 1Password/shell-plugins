package sapbtp

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SAPBTPCLI() schema.Executable {
	return schema.Executable{
		Name:    "SAP BTP CLI",
		Runs:    []string{"btp"},
		DocsURL: sdk.URL("https://help.sap.com/docs/btp/sap-business-technology-platform/account-administration-using-sap-btp-command-line-interface-btp-cli"), // TODO: Replace with actual URL
		NeedsAuth: needsauth.IfAll(
			needsauth.ForCommand("login"),
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
