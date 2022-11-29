package vsce

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "vsce",
		Platform: schema.PlatformInfo{
			Name:     "VS Code Extensions",
			Homepage: sdk.URL("https://www.npmjs.com/package/vsce"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			VSCECLI(),
		},
	}
}
