package nomad

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "nomad",
		Platform: schema.PlatformInfo{
			Name:     "HashiCorp Nomad",
			Homepage: sdk.URL("https://nomad.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			HashiCorpNomadCLI(),
		},
	}
}
