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
				"VAULT_TOKEN":     "jWIJtxke7DkUJ2A9FjJVU9YfmZRF3p04FbsEXAMPLE",
				"VAULT_ADDR":      "https://vault.acme.com",
				"VAULT_NAMESPACE": "default",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token:     "jWIJtxke7DkUJ2A9FjJVU9YfmZRF3p04FbsEXAMPLE",
						fieldname.Address:   "https://vault.acme.com",
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
				fieldname.Token:     "jWIJtxke7DkUJ2A9FjJVU9YfmZRF3p04FbsEXAMPLE",
				fieldname.Address:   "https://vault.acme.com",
				fieldname.Namespace: "default",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"VAULT_TOKEN":     "jWIJtxke7DkUJ2A9FjJVU9YfmZRF3p04FbsEXAMPLE",
					"VAULT_ADDR":      "https://vault.acme.com",
					"VAULT_NAMESPACE": "default",
				},
			},
		},
	})
}
