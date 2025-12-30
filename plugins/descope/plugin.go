package descope

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "descope",
		Platform: schema.PlatformInfo{
			Name:     "Desope",
			Homepage: sdk.URL("https://descope.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			&#34;ManagementKey&#34;(),
		},
		Executables: []schema.Executable{
			DesopeCLI(),
		},
	}
}
