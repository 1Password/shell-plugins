package vault

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "vault",
		Platform: schema.PlatformInfo{
			Name:     "HashiCorp Vault",
			Homepage: sdk.URL("https://www.vaultproject.io"),
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			Executable_vault(),
		},
	}
}
