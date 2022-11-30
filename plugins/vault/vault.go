package vault

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func VaultCLI() schema.Executable {
	return schema.Executable{
		Name:      "Vault CLI",
		Runs:      []string{"vault"},
		DocsURL:   sdk.URL("https://developer.hashicorp.com/vault/docs/commands"),
		NeedsAuth: needsauth.NotForHelpOrVersion(),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
