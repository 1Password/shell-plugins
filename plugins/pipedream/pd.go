package pipedream

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func PipedreamCLI() schema.Executable {
	return schema.Executable{
		Name:    "Pipedream CLI",
		Runs:    []string{"pd"},
		DocsURL: sdk.URL("https://pipedream.com/docs/cli/reference/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
