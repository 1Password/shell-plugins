package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_Okta() schema.Executable {
	return schema.Executable{
		Runs:    []string{"okta"},
		Name:    "Okta CLI",
		DocsURL: sdk.URL("https://cli.okta.com"),
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		NeedsAuth: needsauth.NotForHelpOrVersion(),
	}
}
