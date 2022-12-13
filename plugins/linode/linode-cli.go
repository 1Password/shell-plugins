package linode

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LinodeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Linode CLI",
		Runs:      []string{"linode-cli"},
		DocsURL:   sdk.URL("https://www.linode.com/docs/products/tools/cli/get-started/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
