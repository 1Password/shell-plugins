package anthropic

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ClaudeCodeCLI() schema.Executable {
	return schema.Executable{
		Name:      "Claude Code",
		Runs:      []string{"claude"},
		DocsURL:   sdk.URL("https://docs.anthropic.com/en/docs/agents-and-tools/claude-code/overview"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
