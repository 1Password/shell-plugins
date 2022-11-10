package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OktaCLI() schema.Executable {
	return schema.Executable{
		Runs:    []string{"okta"},
		Name:    "Okta CLI",
		DocsURL: sdk.URL("https://cli.okta.com"),
		UsesCredentials: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
		NeedsAuth: needsauth.NotForHelpOrVersion(),
	}
}
