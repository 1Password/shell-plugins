package vault

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"VAULT_TOKEN":     "abcd123",
				"VAULT_ADDR":      "localhost",
				"VAULT_NAMESPACE": "default",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token:     "abcd123",
						fieldname.Address:   "localhost",
						fieldname.Namespace: "default",
					},
				},
			},
		},
	})
}

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().Provisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[string]string{
				fieldname.Token:     "abcd123",
				fieldname.Address:   "localhost",
				fieldname.Namespace: "default",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"VAULT_TOKEN":     "abcd123",
					"VAULT_ADDR":      "localhost",
					"VAULT_NAMESPACE": "default",
				},
			},
		},
	})
}
