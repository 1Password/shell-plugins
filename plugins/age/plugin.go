package age

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "age",
		Platform: schema.PlatformInfo{
			Name:     "Age",
			Homepage: sdk.URL("https://age-encryption.org/"),
		},
		Credentials: []schema.CredentialType{
			PrivateKey(),
		},
		Executables: []schema.Executable{
			AgeCLI(),
		},
	}
}
