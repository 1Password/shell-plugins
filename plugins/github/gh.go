package github

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_gh() schema.Executable {
	return schema.Executable{
		Runs:    []string{"gh"},
		Name:    "GitHub CLI",
		DocsURL: sdk.URL("https://cli.github.com"),
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		NeedsAuth: needsauth.NotForHelpOrVersion(),
	}
}
