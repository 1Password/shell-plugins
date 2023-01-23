package openai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenAICLI() schema.Executable {
	return schema.Executable{
		Name:      "OpenAI CLI",
		Runs:      []string{"openai"},
		DocsURL:   sdk.URL("https://pypi.org/project/openai/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
