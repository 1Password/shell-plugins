package lacework

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func LaceworkCLI() schema.Executable {
	return schema.Executable{
		Name:      "Lacework CLI",
		Runs:      []string{"lacework"},
		DocsURL:   sdk.URL("https://docs.lacework.com/cli/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
