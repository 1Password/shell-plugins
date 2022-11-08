package github

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "github",
		Platform: schema.PlatformInfo{
			Name:     "GitHub",
			Homepage: sdk.URL("https://github.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			GitHubCLI(),
		},
	}
}
