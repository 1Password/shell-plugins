package sourcegraph

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "sourcegraph",
		Platform: schema.PlatformInfo{
			Name:     "Sourcegraph",
			Homepage: sdk.URL("https://sourcegraph.com"),
		},
		Credentials: []schema.CredentialType{
			AccessToken(),
		},
		Executables: []schema.Executable{
			SourcegraphCLI(),
		},
	}
}
