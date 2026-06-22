package exercism

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "exercism",
		Platform: schema.PlatformInfo{
			Name:     "Exercism",
			Homepage: sdk.URL("https://exercism.org"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			ExercismCLI(),
		},
	}
}
