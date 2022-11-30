package gitlab

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GitLabCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"glab"},
		Name:      "GitLab CLI",
		DocsURL:   sdk.URL("https://glab.readthedocs.io"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
