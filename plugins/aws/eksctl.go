package aws

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func eksctlCLI() schema.Executable {
	return schema.Executable{
		Name:    "eksctl CLI",
		Runs:    []string{"eksctl"},
		DocsURL: sdk.URL("https://eksctl.io/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.AccessKey,
				Provisioner: CLIProvisioner{},
			},
		},
	}
}
