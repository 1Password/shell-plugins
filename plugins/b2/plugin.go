package b2

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "b2",
		Platform: schema.PlatformInfo{
			Name:     "Backblaze B2",
			Homepage: sdk.URL("https://www.backblaze.com/"),
		},
		Credentials: []schema.CredentialType{
			ApplicationKey(),
		},
		Executables: []schema.Executable{
			BackblazeB2CLI(),
		},
	}
}
