package hcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HetznerCloudCLI() schema.Executable {
	return schema.Executable{
		Name:      "Hetzner Cloud CLI",
		Runs:      []string{"hcloud"},
		DocsURL:   sdk.URL("https://github.com/hetznercloud/cli"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
