package ansiblevault

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func AnsibleVaultCLI() schema.Executable {
	return schema.Executable{
		Name:    "Ansible Vault CLI",
		Runs:    []string{"ansible-vault"},
		DocsURL: sdk.URL("https://docs.ansible.com/ansible/latest/cli/ansible-vault.html"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.Credential,
			},
		},
	}
}
