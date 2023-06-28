package render

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "render",
		Platform: schema.PlatformInfo{
			Name:     "Render",
			Homepage: sdk.URL("https://render.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			RenderCLI(),
		},
	}
}
