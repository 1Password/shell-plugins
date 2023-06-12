package pulumi

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

var PulumiSubCommandsNeedingAuth = needsauth.IfAny(
	needsauth.ForCommand("destroy"),
	needsauth.ForCommand("import"),
	needsauth.ForCommand("new"),
	needsauth.ForCommand("org"),
	needsauth.ForCommand("preview"),
	needsauth.ForCommand("refresh"),
	needsauth.ForCommand("stack"),
	needsauth.ForCommand("state"),
	needsauth.ForCommand("up"),
	needsauth.ForCommand("whoami"),
)

func PulumiCLI() schema.Executable {
	return schema.Executable{
		Name:    "Pulumi CLI",
		Runs:    []string{"pulumi"},
		DocsURL: sdk.URL("https://www.pulumi.com/docs/reference/cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name:        credname.PersonalAccessToken,
				Description: "Pulumi state backend configuration (token)",
				Optional:    true,
				NeedsAuth:   PulumiSubCommandsNeedingAuth,
			},
			{
				Name:        BackendOnlyCredentialName,
				Description: "Pulumi state backend configuration (backend)",
				Optional:    true,
				NeedsAuth:   PulumiSubCommandsNeedingAuth,
			},
			{
				Description: "Credentials to use within the Pulumi project",
				SelectFrom: &schema.CredentialSelection{
					ID:                    "project",
					IncludeAllCredentials: true,
					AllowMultiple:         true,
				},
				Optional:  true,
				NeedsAuth: PulumiSubCommandsNeedingAuth,
			},
		},
	}
}
