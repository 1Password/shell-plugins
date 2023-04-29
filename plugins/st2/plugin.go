package st2

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "st2",
		Platform: schema.PlatformInfo{
			Name:     "StackStorm",
			Homepage: sdk.URL("https://stackstorm.com"),
		},
		Credentials: []schema.CredentialType{
			UserPass(),
		},
		Executables: []schema.Executable{
			StackStormCLI(),
		},
	}
}
