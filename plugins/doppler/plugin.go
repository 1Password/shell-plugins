package doppler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "doppler",
		Platform: schema.PlatformInfo{
			Name:     "Doppler",
			Homepage: sdk.URL("https://doppler.com"),
		},
		Credentials: []schema.CredentialType{
			ServiceToken(),
		},
		Executables: []schema.Executable{
			DopplerCLI(),
		},
	}
}
