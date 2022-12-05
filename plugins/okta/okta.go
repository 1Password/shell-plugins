package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OktaCLI() schema.Executable {
	return schema.Executable{
		Name:      "Okta CLI",
		Runs:      []string{"okta"},
		DocsURL:   sdk.URL("https://cli.okta.com"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
