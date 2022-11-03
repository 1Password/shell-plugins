package snyk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_snyk() schema.Executable {
	return schema.Executable{
		Runs:      []string{"snyk"},
		Name:      "Snyk CLI",
		DocsURL:   sdk.URL("https://docs.snyk.io/snyk-cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			APIToken(),
		},
	}
}
