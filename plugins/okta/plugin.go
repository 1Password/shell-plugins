package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "okta",
		Platform: schema.PlatformInfo{
			Name:     "Okta",
			Homepage: sdk.URL("https://www.okta.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			OktaCLI(),
		},
	}
}
