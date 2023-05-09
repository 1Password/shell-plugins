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
			needsauth.NotWhenContainsArgs("login"),
			needsauth.NotWhenContainsArgs("logout"),
		),
		SupportedCredentialAmount: schema.DynamicNumber,
	}
}
