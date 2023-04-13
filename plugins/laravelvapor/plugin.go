package laravelvapor

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "laravelvapor",
		Platform: schema.PlatformInfo{
			Name:     "Laravel Vapor",
			Homepage: sdk.URL("https://vapor.laravel.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
		},
		Executables: []schema.Executable{
			LaravelVaporCLI(),
		},
	}
}
