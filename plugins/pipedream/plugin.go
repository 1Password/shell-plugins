package pipedream

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "pipedream",
		Platform: schema.PlatformInfo{
			Name:     "Pipedream",
			Homepage: sdk.URL("https://pipedream.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			PipedreamCLI(),
		},
	}
}
