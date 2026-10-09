package qodana

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "qodana",
		Platform: schema.PlatformInfo{
			Name:     "Qodana",
			Homepage: sdk.URL("https://qodana.cloud"),
		},
		Credentials: []schema.CredentialType{
			ProjectToken(),
		},
		Executables: []schema.Executable{
			QodanaCLI(),
		},
	}
}
