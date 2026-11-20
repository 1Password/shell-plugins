package kubectl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func KubectlCLI() schema.Executable {
	return schema.Executable{
		Name:      "kubectl",
		Runs:      []string{"kubectl"},
		DocsURL:   sdk.URL("https://kubernetes.io/docs/reference/kubectl/"),
		NeedsAuth: needsAuth,
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credentials,
			},
		},
	}
}
