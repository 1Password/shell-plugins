package evals

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "evals",
		Platform: schema.PlatformInfo{
			Name:     "OpenAI Evals",
			Homepage: sdk.URL("https://github.com/openai/evals"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			OpenAIEvalsCLI(),
		},
	}
}
