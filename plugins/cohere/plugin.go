package cohere

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "cohere",
		Platform: schema.PlatformInfo{
			Name:     "Cohere",
			Homepage: sdk.URL("https://cohere.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			CohereCLI(),
		},
	}
}
