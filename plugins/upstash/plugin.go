package upstash

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "upstash",
		Platform: schema.PlatformInfo{
			Name:     "Upstash",
			Homepage: sdk.URL("https://upstash.com"),
		},
		Credentials: []schema.CredentialType{
			APIKey(),
		},
		Executables: []schema.Executable{
			UpstashCLI(),
		},
	}
}
