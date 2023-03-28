package pulumi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "pulumi",
		Platform: schema.PlatformInfo{
			Name:     "Pulumi",
			Homepage: sdk.URL("https://www.pulumi.com"),
		},
		Credentials: []schema.CredentialType{
			PulumiAccessToken(),
		},
		Executables: []schema.Executable{
			PulumiCLI(),
		},
	}
}
