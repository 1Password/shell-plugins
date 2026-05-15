package todoist

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TodoistCLI() schema.Executable {
	return schema.Executable{
		Name:    "Todoist CLI",
		Runs:    []string{"todoist"},
		DocsURL: sdk.URL("https://github.com/sachaos/todoist"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
