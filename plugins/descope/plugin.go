package descope

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "descope",
		Platform: schema.PlatformInfo{
			Name:     "Decsope",
			Homepage: sdk.URL("https://descope.com"),
		},
		Credentials: []schema.CredentialType{
			ManagementKey(),
		},
		Executables: []schema.Executable{
			DescopeCLI(),
		},
	}
}
