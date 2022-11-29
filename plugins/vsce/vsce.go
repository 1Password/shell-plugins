package vsce

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func VSCECLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"vsce"},
		Name:      "VS Code Extensions CLI",
		DocsURL:   sdk.URL("https://github.com/microsoft/vscode-vsce"),
		NeedsAuth: needsauth.ForCommands([]string{"publish", "verify-pat"}),
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
	}
}
