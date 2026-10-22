package kubectl

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "kubectl",
		Platform: schema.PlatformInfo{
			Name:     "Kubernetes",
			Homepage: sdk.URL("https://kubernetes.io"),
		},
		Credentials: []schema.CredentialType{
			Credentials(),
		},
		Executables: []schema.Executable{
			KubectlCLI(),
		},
	}
}
