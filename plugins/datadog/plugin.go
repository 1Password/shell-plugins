package datadog

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "datadog",
		Platform: schema.PlatformInfo{
			Name:     "Datadog",
			Homepage: sdk.URL("https://datadoghq.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			Executable_dogshell(),
		},
	}
}
