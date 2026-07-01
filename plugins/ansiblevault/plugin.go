package ansiblevault

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "ansiblevault",
		Platform: schema.PlatformInfo{
			Name:     "Ansible Vault",
			Homepage: sdk.URL("https://docs.ansible.com/ansible/latest/cli/ansible-vault.html"),
		},
		Credentials: []schema.CredentialType{
			Password(),
		},
		Executables: []schema.Executable{
			AnsibleVaultCLI(),
		},
	}
}
