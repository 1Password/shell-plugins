package descope

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func DescopeCLI() schema.Executable {
	return schema.Executable{
		Name:    "Descope CLI",
		Runs:    []string{"descope"},
		DocsURL: sdk.URL("https://docs.descope.com/cli/descope"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.ManagementKey,
			},
		},
	}
}
