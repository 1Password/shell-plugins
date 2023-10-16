package vertica

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "vertica",
		Platform: schema.PlatformInfo{
			Name:     "Vertica",
			Homepage: sdk.URL("https://www.vertica.com/"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			VerticaCLI(),
		},
	}
}
