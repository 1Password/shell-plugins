package flyctl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FlyCLI() schema.Executable {
	return schema.Executable{
		Name:      "Fly",
		Runs:      []string{"fly"},
		DocsURL:   sdk.URL("https://fly.io/docs/flyctl/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
