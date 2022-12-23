package gitlab

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GitLabCLI() schema.Executable {
	return schema.Executable{
		Name:    "GitLab CLI",
		Runs:    []string{"glab"},
		DocsURL: sdk.URL("https://glab.readthedocs.io"),
		NeedsAuth: needsauth.For(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
