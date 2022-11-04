package okta

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "okta",
		Platform: schema.PlatformInfo{
			Name:     "Okta CLI",
			Homepage: sdk.URL("https://www.okta.com"),
			Logo:     sdk.URL("https://www.okta.com/themes/custom/okta_www_theme/images/logo.svg?v2"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			Executable_Okta(),
		},
	}
}
