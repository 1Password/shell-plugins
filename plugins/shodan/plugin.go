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
			Homepage: sdk.URL("https://shodan.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ShodanCLI(),
		},
	}
}
