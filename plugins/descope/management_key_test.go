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
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Key: "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DESCOPE_KEY": "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
				},
			},
		},
	})
}

func TestManagementKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, ManagementKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"DESCOPE_KEY": "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
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
