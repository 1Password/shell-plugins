package axiom

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "axiom",
		Platform: schema.PlatformInfo{
			Name:     "Axiom",
			Homepage: sdk.URL("https://axiom.co"),
		},
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
		Executables: []schema.Executable{
			AxiomCLI(),
		},
	}
}
