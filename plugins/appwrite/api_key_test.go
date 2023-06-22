package appwrite

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
				fieldname.APIKey: "5gywx05hyh46gexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"APPWRITE_API_KEY": "5gywx05hyh46gexample",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"APPWRITE_API_KEY": "5gywx05hyh46gexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "5gywx05hyh46gexample",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in appwrite/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "5gywx05hyh46gexample",
			// 		},
			// 	},
			},
		},
	})
}
