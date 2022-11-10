package github

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GitHubCLI() schema.Executable {
	return schema.Executable{
		Runs:    []string{"gh"},
		Name:    "GitHub CLI",
		DocsURL: sdk.URL("https://cli.github.com"),
		UsesCredentials: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
		NeedsAuth: needsauth.NotForHelpOrVersion(),
	}
}
