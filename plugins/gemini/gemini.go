package gemini

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func GoogleGeminiCLI() schema.Executable {
	return schema.Executable{
		Name:    "Google Gemini CLI",
		Runs:    []string{"gemini"},
		DocsURL: sdk.URL("https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/index.md"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
