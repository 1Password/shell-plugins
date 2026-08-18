package grok

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "grok",
		Platform: schema.PlatformInfo{
			Name:     "xAI Grok",
			Homepage: sdk.URL("https://x.ai"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			GrokCLI(),
		},
	}
}
