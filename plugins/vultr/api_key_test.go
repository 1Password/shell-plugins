package vultr

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "4T7FVVVS23ZMKHLV990NL61F4O6L9EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"VULTR_API_KEY": "4T7FVVVS23ZMKHLV990NL61F4O6L9EXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"VULTR_API_KEY": "4T7FVVVS23ZMKHLV990NL61F4O6L9EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "4T7FVVVS23ZMKHLV990NL61F4O6L9EXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.vultr-cli.yaml": plugintest.LoadFixture(t, "vultr-cli.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "4T7FVVVS23ZMKHLV990NL61F4O6L9EXAMPLE",
					},
				},
			},
		},
	})
}

func TestAPIKeyNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, VultrCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no for --config": {
			Args:              []string{"--config"},
			ExpectedNeedsAuth: false,
		},
	})
}
