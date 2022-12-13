package lacework

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "lacework",
		Platform: schema.PlatformInfo{
			Name:     "Lacework",
			Homepage: sdk.URL("https://www.lacework.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			LaceworkCLI(),
		},
	}
}
