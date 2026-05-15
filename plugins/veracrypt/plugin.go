package veracrypt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "veracrypt",
		Platform: schema.PlatformInfo{
			Name:     "VeraCrypt",
			Homepage: sdk.URL("https://www.veracrypt.fr"),
		},
		Credentials: []schema.CredentialType{
			VolumePassword(),
		},
		Executables: []schema.Executable{
			VeraCryptCLI(),
		},
	}
}