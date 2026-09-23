package octopus

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OctopusDeployCLI() schema.Executable {
	return schema.Executable{
		Name:    "Octopus CLI",
		Runs:    []string{"octopus"},
		DocsURL: sdk.URL("https://octopus.com/docs/octopus-rest-api/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
			needsauth.NotWhenContainsArgs("config"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
