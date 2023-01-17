package sourcegraph

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SourcegraphCLI() schema.Executable {
	return schema.Executable{
		Name:      "Sourcegraph CLI",
		Runs:      []string{"src"},
		DocsURL:   sdk.URL("https://docs.sourcegraph.com/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
