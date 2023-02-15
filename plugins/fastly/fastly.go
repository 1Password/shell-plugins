package fastly

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FastlyCLI() schema.Executable {
	return schema.Executable{
		Name:    "Fastly CLI",
		Runs:    []string{"fastly"},
		DocsURL: sdk.URL("https://developer.fastly.com/reference/cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("-t"),
			needsauth.NotWhenContainsArgs("--token"),
			needsauth.NotWhenContainsArgs("profile"),
			needsauth.NotWhenContainsArgs("config"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
