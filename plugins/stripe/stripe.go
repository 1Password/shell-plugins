package stripe

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func StripeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Stripe CLI",
		Runs:      []string{"stripe"},
		DocsURL:   sdk.URL("https://stripe.com/docs/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
