package kaggle

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func KaggleCLI() schema.Executable {
	return schema.Executable{
		Name:    "Kaggle CLI",
		Runs:    []string{"kaggle"},
		DocsURL: sdk.URL("https://github.com/Kaggle/kaggle-api"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.APIToken,
			},
		},
	}
}
