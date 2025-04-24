package openai

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "openai",
		Platform: schema.PlatformInfo{
			Name:     "OpenAI",
			Homepage: sdk.URL("https://openai.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			OpenAICLI(),
			OpenAIEvalsCLI(),
			OpenAIEvalSetCLI(),
			OpenAICodex(),
		},
	}
}
