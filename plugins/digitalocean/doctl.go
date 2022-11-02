package digitalocean

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func Executable_doctl() schema.Executable {
	return schema.Executable{
		Runs:      []string{"doctl"},
		Name:      "DigitalOcean CLI",
		DocsURL:   sdk.URL("https://docs.digitalocean.com/reference/doctl"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Credentials: []schema.CredentialType{
			PersonalAccessToken(),
		},
	}
}
