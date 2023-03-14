package zendesk

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ZENDESK_TOKEN": "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"ZENDESK_TOKEN": "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in zendesk/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
