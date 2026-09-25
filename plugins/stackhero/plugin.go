package stackhero

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "stackhero",
		Platform: schema.PlatformInfo{
			Name:     "Stackhero",
			Homepage: sdk.URL("https://www.stackhero.io/"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			StackheroCLI(),
		},
	}
}
