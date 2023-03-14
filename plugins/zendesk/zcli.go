package zendesk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ZendeskCLI() schema.Executable {
	return schema.Executable{
		Name:    "Zendesk CLI",
		Runs:    []string{"zcli"},
		DocsURL: sdk.URL("https://developer.zendesk.com/documentation/apps/getting-started/using-zcli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
