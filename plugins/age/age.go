package age

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AgeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Age CLI",
		Runs:    []string{"age"},
		DocsURL: sdk.URL("https://age-encryption.org"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotForExactArgs("encrypt", "recipient", "armor"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
