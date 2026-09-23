package pypi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func FlitCLI() schema.Executable {
	return schema.Executable{
		Name:    "Flit",
		Runs:    []string{"flit"},
		DocsURL: sdk.URL("https://flit.pypa.io"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.ForCommand("publish"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.APIToken,
				Provisioner: PyPIToolProvisioner("FLIT_USERNAME", "FLIT_PASSWORD"),
			},
		},
	}
}
