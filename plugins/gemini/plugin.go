package gemini

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "gemini",
		Platform: schema.PlatformInfo{
			Name:     "Google Gemini",
			Homepage: sdk.URL("https://gemini.google.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			GoogleGeminiCLI(),
		},
	}
}
