package spacelift

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "spacelift",
		Platform: schema.PlatformInfo{
			Name:     "Spacelift",
			Homepage: sdk.URL("https://spacelift.io"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			SpaceliftCLI(),
		},
	}
}
