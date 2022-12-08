package tugboat

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "tugboat",
		Platform: schema.PlatformInfo{
			Name:     "Tugboat",
			Homepage: sdk.URL("https://tugboatqa.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			TugboatCLI(),
		},
	}
}
