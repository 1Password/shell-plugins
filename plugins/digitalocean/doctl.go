package digitalocean

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DigitalOceanCLI() schema.Executable {
	return schema.Executable{
		Runs:      []string{"doctl"},
		Name:      "DigitalOcean CLI",
		DocsURL:   sdk.URL("https://docs.digitalocean.com/reference/doctl"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
