package linode

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "linode",
		Platform: schema.PlatformInfo{
			Name:     "Linode",
			Homepage: sdk.URL("https://linode.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			LinodeCLI(),
		},
	}
}
