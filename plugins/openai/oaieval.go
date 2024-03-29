package openai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func OpenAIEvalsCLI() schema.Executable {
	return schema.Executable{
		Name:    "OpenAI Evals CLI",
		Runs:    []string{"oaieval"},
		DocsURL: sdk.URL("https://github.com/openai/evals/blob/main/docs/run-evals.md"),
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
