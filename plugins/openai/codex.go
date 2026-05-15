package openai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenAICodex() schema.Executable {
	return schema.Executable{
		Name:      "OpenAI Codex",
		Runs:      []string{"codex"},
		DocsURL:   sdk.URL("https://help.openai.com/en/articles/11096431-openai-codex-cli-getting-started"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
