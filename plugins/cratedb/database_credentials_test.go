package cratedb

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.: "yqch",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CRATEDB": "yqch",
				},
			},
		},
	})
}

func TestDatabaseCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"CRATEPW": "eog-l4ogPoIO4kX8ICHC*6kP",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "eog-l4ogPoIO4kX8ICHC*6kP",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in cratedb/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "yqch",
			// 		},
			// 	},
			},
		},
	})
}
