package descope

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestManagementKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ManagementKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ 
				fieldname.ProjectID: "P37a8ZDXSo0XPrH3maK9PvF5IKrNEXAMPLE",
				fieldname.ManagementKey: "K37aAMWGYC8trD7MXBp2P9y22kBsi0qnNMRQQVLSK7YxbT9taEHnTpfrJTb2Qozm9Yd4sAeEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DESCOPE_PROJECT_ID": "P37a8ZDXSo0XPrH3maK9PvF5IKrNEXAMPLE",
					"DESCOPE_MANAGEMENT_KEY": "K37aAMWGYC8trD7MXBp2P9y22kBsi0qnNMRQQVLSK7YxbT9taEHnTpfrJTb2Qozm9Yd4sAeEXAMPLE",
				},
			},
		},
	})
}

func TestManagementKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, ManagementKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ 
				"DESCOPE_PROJECT_ID": "P37a8ZDXSo0XPrH3maK9PvF5IKrNEXAMPLE",
				"DESCOPE_MANAGEMENT_KEY": "K37aAMWGYC8trD7MXBp2P9y22kBsi0qnNMRQQVLSK7YxbT9taEHnTpfrJTb2Qozm9Yd4sAeEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ProjectID: "P37a8ZDXSo0XPrH3maK9PvF5IKrNEXAMPLE",
				fieldname.ManagementKey: "K37aAMWGYC8trD7MXBp2P9y22kBsi0qnNMRQQVLSK7YxbT9taEHnTpfrJTb2Qozm9Yd4sAeEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in descope/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
				// 		},
				// 	},
			},
		},
	})
}
