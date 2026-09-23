package gemini

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "gemini",
		Platform: schema.PlatformInfo{
			Name:     "Google Gemini CLI",
			Homepage: sdk.URL("https://geminicli.com/"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			GoogleGeminiCLI(),
		},
	}
}
