package pypi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func TwineCLI() schema.Executable {
	return schema.Executable{
		Name:    "Twine",
		Runs:    []string{"twine"},
		DocsURL: sdk.URL("https://twine.readthedocs.io"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.ForCommand("upload"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.APIToken,
				Provisioner: PyPIToolProvisioner("TWINE_USERNAME", "TWINE_PASSWORD"),
			},
		},
	}
}
