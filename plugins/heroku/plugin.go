package heroku

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "heroku",
		Platform: schema.PlatformInfo{
			Name:     "Heroku",
			Homepage: sdk.URL("https://heroku.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			Executable_heroku(),
		},
	}
}
