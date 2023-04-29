package zendesk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "zendesk",
		Platform: schema.PlatformInfo{
			Name:     "Zendesk",
			Homepage: sdk.URL("https://www.zendesk.com/"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			ZendeskCLI(),
		},
	}
}
