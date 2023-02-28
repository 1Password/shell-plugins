package ngrok

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ngrok",
		Platform: schema.PlatformInfo{
			Name:     "ngrok",
			Homepage: sdk.URL("https://ngrok.com"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
			CredentialsAPI(),
		},
		Executables: []schema.Executable{
			ngrokCLI(),
		},
	}
}
