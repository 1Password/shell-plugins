package tugboat

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TugboatCLI() schema.Executable {
	return schema.Executable{
		Name:      "Tugboat CLI",
		Runs:      []string{"tugboat"},
		DocsURL:   sdk.URL("https://docs.tugboatqa.com/tugboat-cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AccessToken,
			},
		},
	}
}
