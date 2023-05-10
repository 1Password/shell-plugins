package terraform

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credselect"
)

func TerraformCLI() schema.Executable {
	return schema.Executable{
		Name:    "Terraform CLI",
		Runs:    []string{"terraform"},
		DocsURL: sdk.URL("https://developer.hashicorp.com/terraform/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Description: "Credentials to use within the Terraform project",
				SelectFrom: credselect.New("project",
					// Allow any credential to be provisioned to Terraform
					credselect.Any(),
					// Allow multiple credentials to be provisioned at the same time
					credselect.AllowMultiple(),
				),
				NeedsAuth: needsauth.IfAny(
					needsauth.ForCommand("refresh"),
					needsauth.ForCommand("plan"),
					needsauth.ForCommand("apply"),
				),
			},
		},
	}
}
