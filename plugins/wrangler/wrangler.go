package wrangler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func CloudflareWorkersCLI() schema.Executable {
	return schema.Executable{
		Name:      "Cloudflare Workers CLI",
		Runs:      []string{"wrangler"},
		DocsURL:   sdk.URL("https://wrangler.com/docs/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
