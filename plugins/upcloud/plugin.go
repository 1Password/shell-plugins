package upcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "upcloud",
		Platform: schema.PlatformInfo{
			Name:     "UpCloud",
			Homepage: sdk.URL("https://upcloud.com"),
		},
		Credentials: []schema.CredentialType{
			UserLogin(),
		},
		Executables: []schema.Executable{
			UpCloudCLI(),
		},
	}
}
