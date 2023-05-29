package vercel

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "vercel",
		Platform: schema.PlatformInfo{
			Name:     "Vercel",
			Homepage: sdk.URL("https://vercel.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			VercelCLI(),
		},
	}
}
