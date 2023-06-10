package spacelift

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func SpaceliftCLI() schema.Executable {
	return schema.Executable{
		Name:      "Spacelift CLI",
		Runs:      []string{"spacectl"},
		DocsURL:   sdk.URL("https://github.com/spacelift-io/spacectl"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIKey,
			},
		},
	}
}
