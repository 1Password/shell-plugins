package nomad

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HashiCorpNomadCLI() schema.Executable {
	return schema.Executable{
		Name:      "HashiCorp Nomad CLI",
		Runs:      []string{"nomad"},
		DocsURL:   sdk.URL("https://developer.hashicorp.com/nomad/docs/commands"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
