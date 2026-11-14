package bundler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "bundler",
		Platform: schema.PlatformInfo{
			Name:     "Ruby Bundler",
			Homepage: sdk.URL("https://bundler.io"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			BundleCLI(),
		},
	}
}
