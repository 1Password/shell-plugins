package argocd

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func ArgocdCLI() schema.Executable {
	return schema.Executable{
		Name:      "Argo CD CLI",
		Runs:      []string{"argocd"},
		DocsURL:   sdk.URL("https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd/"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
