package circleci

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "circleci",
		Platform: schema.PlatformInfo{
			Name:     "CircleCI",
			Homepage: sdk.URL("https://circleci.com"),
		},
		Credentials: []schema.CredentialType{
			PersonalAPIToken(),
		},
		Executables: []schema.Executable{
			Executable_circleci(),
		},
	}
}
