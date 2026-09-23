package copilot

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CopilotCLI() schema.Executable {
	return schema.Executable{
		Name:    "GitHub Copilot CLI",
		Runs:    []string{"copilot"},
		DocsURL: sdk.URL("https://docs.github.com/en/copilot/reference/copilot-cli-reference/cli-command-reference"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("login"),
			needsauth.NotForExactArgs("version"),
			needsauth.NotWhenContainsArgs("completion"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
