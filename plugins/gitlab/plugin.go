package gitlab

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "gitlab",
		Platform: schema.PlatformInfo{
			Name:     "GitLab",
			Homepage: sdk.URL("https://gitlab.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			GitLabCLI(),
		},
	}
}
