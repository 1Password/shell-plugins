package terragrunt

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func TerragruntCLI() schema.Executable {
	return schema.Executable{
		Name:    "Terragrunt CLI",
		Runs:    []string{"terragrunt"},
		DocsURL: sdk.URL("https://terragrunt.gruntwork.io/docs/reference/cli-options/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Description: "Credentials to use within the Terragrunt project",
				SelectFrom: &schema.CredentialSelection{
					ID:                    "project",
					IncludeAllCredentials: true,
					AllowMultiple:         true,
				},
				Optional: true,
				NeedsAuth: needsauth.IfAny(
					needsauth.ForCommand("refresh"),
					needsauth.ForCommand("init"),
					needsauth.ForCommand("state", "list"),
					needsauth.ForCommand("plan"),
					needsauth.ForCommand("apply"),
					needsauth.ForCommand("destroy"),
					needsauth.ForCommand("import"),
					needsauth.ForCommand("test"),
					needsauth.ForCommand("output-module-groups"),
					needsauth.ForCommand("output-module-groups", "apply"),
					needsauth.ForCommand("output-module-groups", "destroy"),
					needsauth.ForCommand("run-all", "plan"),
					needsauth.ForCommand("run-all", "apply"),
					needsauth.ForCommand("run-all", "output"),
					needsauth.ForCommand("run-all", "destroy"),
					needsauth.ForCommand("run-all", "validate"),
				),
			},
		},
	}
}
