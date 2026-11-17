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
		DocsURL: sdk.URL("https://geminicli.com/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
