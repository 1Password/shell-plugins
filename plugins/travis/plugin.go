package travis

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "travis",
		Platform: schema.PlatformInfo{
			Name:     "Travis CI",
			Homepage: sdk.URL("https://www.travis-ci.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			TravisCICLI(),
		},
	}
}
