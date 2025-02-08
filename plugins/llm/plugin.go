package llm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "llm",
		Platform: schema.PlatformInfo{
			Name:     "LLM",
			Homepage: sdk.URL("https://llm.datasette.io/"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			LLMCLI(),
		},
	}
}
