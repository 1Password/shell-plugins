package azure

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AzureCLI() schema.Executable {
	return schema.Executable{
		Name:    "Azure CLI",
		Runs:    []string{"az"},
		DocsURL: sdk.URL("https://learn.microsoft.com/en-us/cli/azure/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.ServicePrincipal,
				Provisioner: CLIProvisioner{},
			},
		},
	}
}
