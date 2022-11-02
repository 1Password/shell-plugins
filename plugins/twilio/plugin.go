package twilio

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "twilio",
		Platform: schema.PlatformInfo{
			Name:     "Twilio",
			Homepage: sdk.URL("https://twilio.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			Executable_twilio(),
		},
	}
}
