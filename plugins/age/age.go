package age

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AgeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Age",
		Runs:    []string{"age"},
		DocsURL: sdk.URL("https://htmlpreview.github.io/?https://github.com/FiloSottile/age/blob/main/doc/age.1.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
