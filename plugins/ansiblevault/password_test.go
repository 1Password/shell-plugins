package ansiblevault

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPasswordProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Password().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "Sl9WhfECfUn43o2qEcN2GnAw2EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"/tmp/.ansible-vault": {
						Contents: []byte(plugintest.LoadFixture(t, ".ansible-vault")),
					},
				},
				Environment: map[string]string{
					"ANSIBLE_VAULT_PASSWORD_FILE": "/tmp/.ansible-vault",
				},
			},
		},
	})
}

func TestPasswordImporter(t *testing.T) {
	plugintest.TestImporter(t, Password().Importer, map[string]plugintest.ImportCase{
		"default": {
			Files: map[string]string{
				"~/.ansible-vault": plugintest.LoadFixture(t, ".ansible-vault"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "Sl9WhfECfUn43o2qEcN2GnAw2EXAMPLE",
					},
				},
			},
		},
	})
}
