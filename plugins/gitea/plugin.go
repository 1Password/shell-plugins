package gitea

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "gitea",
		Platform: schema.PlatformInfo{
			Name:     "Gitea",
			Homepage: sdk.URL("https://gitea.io/"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			giteaCLI(),
		},
	}
}
