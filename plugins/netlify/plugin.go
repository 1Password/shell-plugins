package netlify

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "netlify",
		Platform: schema.PlatformInfo{
			Name:     "Netlify",
			Homepage: sdk.URL("https://app.netlify.com/"),
		},
		Credentials: []schema.CredentialType{
			PersonalAPIToken(),
		},
		Executables: []schema.Executable{
			NetlifyCLI(),
		},
	}
}
