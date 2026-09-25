package nimble

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "nimble",
		Platform: schema.PlatformInfo{
			Name:     "Nimble",
			Homepage: sdk.URL("https://nimble.directory/"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			nimCLI(),
		},
	}
}
