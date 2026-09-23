package copilot

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "copilot",
		Platform: schema.PlatformInfo{
			Name:     "GitHub Copilot",
			Homepage: sdk.URL("https://github.com/features/copilot/cli"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			CopilotCLI(),
		},
	}
}
