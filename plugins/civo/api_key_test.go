package civo

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.APIKey: "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CIVO_API_KEY": "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"CIVO_API_KEY": "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in civo/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
