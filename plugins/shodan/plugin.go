package shodan

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "shodan",
		Platform: schema.PlatformInfo{
			Name:     "Shodan",
			Homepage: sdk.URL("https://www.shodan.io"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ShodanCLI(),
		},
	}
}
