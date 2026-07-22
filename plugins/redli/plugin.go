package redli

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "redli",
		Platform: schema.PlatformInfo{
			Name:     "Redli",
			Homepage: sdk.URL("https://github.com/IBM-Cloud/redli"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			RedliCLI(),
		},
	}
}
