package sops

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func HelmCLI() schema.Executable {
	return schema.Executable{
		Name:    "Helm with SOPS Secrets",
		Runs:    []string{"helm"},
		DocsURL: sdk.URL("https://github.com/jkroepke/helm-secrets"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			// Only authenticate when using helm-secrets plugin
			needsauth.ForCommand("secrets"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.SecretKey,
			},
		},
	}
}
