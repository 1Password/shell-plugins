package expo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "expo",
		Platform: schema.PlatformInfo{
			Name:     "Expo",
			Homepage: sdk.URL("https://expo.dev"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			ExpoCLI(),
			EASCLI(),
		},
	}
}
