package wrangler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "wrangler",
		Platform: schema.PlatformInfo{
			Name:     "Cloudflare Workers",
			Homepage: sdk.URL("https://workers.cloudflare.com"),
		},
		Credentials: []schema.CredentialType{
			APIToken(),
			APIKey(),
		},
		Executables: []schema.Executable{
			CloudflareWorkersCLI(),
		},
	}
}
