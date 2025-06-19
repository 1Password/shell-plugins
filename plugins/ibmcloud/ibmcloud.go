package ibmcloud

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func IBMCloudCLI() schema.Executable {
	return schema.Executable{
		Name:    "IBM Cloud CLI",
		Runs:    []string{"ibmcloud"},
		DocsURL: sdk.URL("https://ibmcloud.com/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.APIKey,
				Provisioner: CLIProvisioner{},
			},
		},
	}
}
