package terraform

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/credselect"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func TerraformCLI() schema.Executable {
	return schema.Executable{
		Name:      "Terraform CLI",
		Runs:      []string{"terraform"},
		DocsURL:   sdk.URL("https://developer.hashicorp.com/terraform/docs"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Select:        credselect.Any,
				AllowMultiple: true,
				NeedsAuth: needsauth.ForCommands(
					[]string{"plan"},
					[]string{"apply"},
					[]string{"validate"},
					[]string{"destroy"},
					[]string{"import"},
					[]string{"refresh"},
				),
			},
		},
	}
}
