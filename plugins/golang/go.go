package golang

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/credselect"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func GoCLI() schema.Executable {
	return schema.Executable{
		Name:      "Go CLI",
		Runs:      []string{"go"},
		DocsURL:   sdk.URL("https://pkg.go.dev/cmd/go"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Description:   "Authenticate to private module repository.",
				Select:        credselect.Any,
				AllowMultiple: true,
				NeedsAuth: needsauth.ForCommands(
					[]string{"mod", "tidy"},
					[]string{"mod", "download"},
					// TODO: Add more commands that require authentication for pulling in private modules
				),
			},
			{
				Description:   "Provision run command with secrets.",
				Select:        credselect.Any,
				AllowMultiple: true,
				NeedsAuth: needsauth.ForCommands(
					[]string{"run"},
				),
			},
			{
				Description:   "Provision test command with secrets.",
				Select:        credselect.Any,
				AllowMultiple: true,
				NeedsAuth: needsauth.ForCommands(
					[]string{"test"},
				),
			},
		},
	}
}
