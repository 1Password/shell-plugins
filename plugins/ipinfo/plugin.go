package ipinfo

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ipinfo",
		Platform: schema.PlatformInfo{
			Name:     "IPinfo.io",
			Homepage: sdk.URL("https://ipinfo.io"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			IPinfoCLI(),
		},
	}
}
