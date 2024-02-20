package azure

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "azure",
		Platform: schema.PlatformInfo{
			Name:     "Azure",
			Homepage: sdk.URL("https://azure.com"),
		},
		Credentials: []schema.CredentialType{
			ServicePrincipal(),
		},
		Executables: []schema.Executable{
			AzureCLI(),
		},
	}
}
