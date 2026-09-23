package octopus

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
				fieldname.APIKey: "API-PFIPLPKUTCYAWJXUREGWTUYNYEEXAMPLE",
				fieldname.URL:    "https://my.octopus.app",
				fieldname.Space:  "Default",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OCTOPUS_API_KEY": "API-PFIPLPKUTCYAWJXUREGWTUYNYEEXAMPLE",
					"OCTOPUS_URL":     "https://my.octopus.app",
					"OCTOPUS_SPACE":   "Default",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OCTOPUS_API_KEY": "API-PFIPLPKUTCYAWJXUREGWTUYNYEEXAMPLE",
				"OCTOPUS_URL":     "https://my.octopus.app",
				"OCTOPUS_SPACE":   "Default",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "API-PFIPLPKUTCYAWJXUREGWTUYNYEEXAMPLE",
						fieldname.URL:    "https://my.octopus.app",
						fieldname.Space:  "Default",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/octopus/cli_config.json": plugintest.LoadFixture(t, "cli_config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "API-PFIPLPKUTCYAWJXUREGWTUYNYEEXAMPLE",
						fieldname.URL:    "https://my.octopus.app",
						fieldname.Space:  "Default",
					},
				},
			},
		},
	})
}
