package laravelforge

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "laravelforge",
		Platform: schema.PlatformInfo{
			Name:     "Laravel Forge",
			Homepage: sdk.URL("https://forge.laravel.com/"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			LaravelForgeCLI(),
		},
	}
}
