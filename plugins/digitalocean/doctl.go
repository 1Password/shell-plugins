package digitalocean

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DigitalOceanCLI() schema.Executable {
	return schema.Executable{
		Name:      "DigitalOcean CLI",
		Runs:      []string{"doctl"},
		DocsURL:   sdk.URL("https://docs.digitalocean.com/reference/doctl"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.PersonalAccessToken,
			},
		},
	}
}
