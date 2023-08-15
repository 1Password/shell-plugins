package crowdin

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "crowdin",
		Platform: schema.PlatformInfo{
			Name:     "Crowdin",
			Homepage: sdk.URL("https://crowdin.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			CrowdinCLI(),
		},
	}
}
