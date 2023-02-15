package confluent

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "confluent",
		Platform: schema.PlatformInfo{
			Name:     "Confluent",
			Homepage: sdk.URL("https://confluent.io"),
		},
		Credentials: []schema.CredentialType{
			CloudCredentials(),
		},
		Executables: []schema.Executable{
			ConfluentCLI(),
		},
	}
}
