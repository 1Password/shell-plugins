package localstack

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "localstack",
		Platform: schema.PlatformInfo{
			Name:     "LocalStack",
			Homepage: sdk.URL("https://localstack.cloud"),
		},
		// TODO: LocalStack accepts both auth token and api key. When multiple
		// credentials types are supported, update this list to include both
		// options.
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			LocalStackCLI(),
		},
	}
}
