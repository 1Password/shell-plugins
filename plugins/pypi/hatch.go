package pypi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HatchCLI() schema.Executable {
	return schema.Executable{
		Name:    "Hatch",
		Runs:    []string{"hatch"},
		DocsURL: sdk.URL("https://hatch.pypa.io"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.ForCommand("publish"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.APIToken,
				Provisioner: PyPIToolProvisioner("HATCH_INDEX_USER", "HATCH_INDEX_AUTH"),
			},
		},
	}
}
