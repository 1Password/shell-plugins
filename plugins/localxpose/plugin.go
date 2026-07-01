package localxpose

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "localxpose",
		Platform: schema.PlatformInfo{
			Name:     "LocalXpose",
			Homepage: sdk.URL("https://localxpose.io"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			LocalXposeCLI(),
		},
	}
}
