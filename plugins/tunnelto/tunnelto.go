package tunnelto

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func tunneltodevCLI() schema.Executable {
	return schema.Executable{
		Name:    "tunnelto.dev CLI",
		Runs:    []string{"tunnelto"},
		DocsURL: sdk.URL("https://github.com/agrinman/tunnelto"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("-k"),
			needsauth.NotWhenContainsArgs("--key"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
