package twilio

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TwilioCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"twilio"},
		Name:      "Twilio CLI",
		DocsURL:   sdk.URL("https://twilio.com/docs/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
