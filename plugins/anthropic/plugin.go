package anthropic

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "anthropic",
		Platform: schema.PlatformInfo{
			Name:     "Anthropic",
			Homepage: sdk.URL("https://anthropic.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ClaudeCodeCLI(),
		},
	}
}
