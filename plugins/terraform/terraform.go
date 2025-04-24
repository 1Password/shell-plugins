package terraform

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
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
				SelectFrom: &schema.CredentialSelection{
					ID:                    "project",
					IncludeAllCredentials: true,
					AllowMultiple:         true,
				},
				Optional: true,
				NeedsAuth: needsauth.IfAny(
					needsauth.ForCommand("refresh"),
					needsauth.ForCommand("init"),
					needsauth.ForCommand("state"),
					needsauth.ForCommand("plan"),
					needsauth.ForCommand("apply"),
					needsauth.ForCommand("destroy"),
					needsauth.ForCommand("import"),
					needsauth.ForCommand("test"),
				),
			},
		},
	}
}
