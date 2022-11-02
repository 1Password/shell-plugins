package gitlab

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_glab() schema.Executable {
	return schema.Executable{
		Runs:      []string{"glab"},
		Name:      "GitLab CLI",
		DocsURL:   sdk.URL("https://glab.readthedocs.io"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
	}
}
