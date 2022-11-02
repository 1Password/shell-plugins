package example

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_example() schema.Executable {
	return schema.Executable{
		Runs:      []string{"example"},
		Name:      "Example CLI",
		DocsURL:   sdk.URL("http://example.com/docs/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			APIToken(),
		},
	}
}
