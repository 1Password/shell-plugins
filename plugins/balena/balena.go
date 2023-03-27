package balena

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func BalenaCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"balena"},
		Name:      "Balena CLI",
		DocsURL:   sdk.URL("https://www.balena.io/docs/reference/balena-cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			APIKey(),
		},
	}
}
