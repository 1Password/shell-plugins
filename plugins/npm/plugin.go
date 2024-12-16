package npm

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "npm",
		Platform: schema.PlatformInfo{
			Name:     "NPM",
			Homepage: sdk.URL("https://npmjs.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			NPMCLI(),
		},
	}
}
