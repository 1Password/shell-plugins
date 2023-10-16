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
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			LocalStackCLI(),
		},
	}
}
