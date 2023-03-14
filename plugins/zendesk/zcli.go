package zendesk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ZendeskCLI() schema.Executable {
	return schema.Executable{
		Name:      "Zendesk CLI", // TODO: Check if this is correct
		Runs:      []string{"zcli"},
		DocsURL:   sdk.URL("https://zendesk.com/docs/cli"), // TODO: Replace with actual URL
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
