package prefect

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "prefect",
		Platform: schema.PlatformInfo{
			Name:     "Prefect",
			Homepage: sdk.URL("https://app.prefect.cloud/"),
		},
		Credentials: []schema.CredentialType{
			AccessKey(),
		},
		Executables: []schema.Executable{
			PrefectCLI(),
		},
	}
}
