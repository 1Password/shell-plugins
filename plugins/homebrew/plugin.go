package homebrew

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "homebrew",
		Platform: schema.PlatformInfo{
			Name:     "Homebrew",
			Homepage: sdk.URL("https://brew.sh/"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			HomebrewCLI(),
		},
	}
}
