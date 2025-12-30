package descope

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func Test&#34;ManagementKey&#34;Provisioner(t *testing.T) {
	plugintest.TestProvisioner(t, &#34;ManagementKey&#34;().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Key&#34;: "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DESCOPE_KEY&#34;": "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
				},
			},
		},
	})
}

func Test&#34;ManagementKey&#34;Importer(t *testing.T) {
	plugintest.TestImporter(t, &#34;ManagementKey&#34;().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"DESCOPE_KEY&#34;": "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key&#34;: "Bn3JbRL6odvVBF7ZtzC3lEuQbri2AFtwUpn9tNUj9lE7lbCuk3cYCaPOUZZsdh3bEXAMPLE",
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
