package tofu

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func TofuCLI() schema.Executable {
	return schema.Executable{
		Name:    "OpenTofu",
		Runs:    []string{"tofu"},
		DocsURL: sdk.URL("https://opentofu.org/docs/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Description: "Credentials to use within the OpenTofu project",
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
