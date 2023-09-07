package zapier

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "zapier",
		Platform: schema.PlatformInfo{
			Name:     "Zapier",
			Homepage: sdk.URL("https://zapier.com"),
		},
		Credentials: []schema.CredentialType{
			DeployKey(),
		},
		Executables: []schema.Executable{
			ZapierCLI(),
		},
	}
}
