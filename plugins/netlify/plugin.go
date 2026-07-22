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
			Homepage: sdk.URL("https://netlify.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			NetlifyCLI(),
		},
	}
}
