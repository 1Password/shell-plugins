package argocd

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "argocd",
		Platform: schema.PlatformInfo{
			Name:     "Argo CD",
			Homepage: sdk.URL("https://argo-cd.readthedocs.io/en/stable/"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			ArgocdCLI(),
		},
	}
}
