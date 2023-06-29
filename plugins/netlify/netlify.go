package netlify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func NetlifyCLI() schema.Executable {
	return schema.Executable{
		Name:    "Netlify CLI",
		Runs:    []string{"netlify"},
		DocsURL: sdk.URL("https://docs.netlify.com/cli/get-started/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAPIToken,
			},
		},
	}
}
