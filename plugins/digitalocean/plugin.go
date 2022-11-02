package digitalocean

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "digitalocean",
		Platform: schema.PlatformInfo{
			Name:     "DigitalOcean",
			Homepage: sdk.URL("https://www.digitalocean.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			Executable_doctl(),
		},
	}
}
