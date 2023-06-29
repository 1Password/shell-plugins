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
			ItemFields: map[sdk.FieldName]string{ 
				fieldname.Password: "1<34&f0rg3t@me",
				fieldname.Username: "admin",
				fieldname.Host: "https://love.aks1.eastus2.azure.cratedb.net:4200",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CRATEPW": "1<34&f0rg3t@me",
				},
				CommandLine: []string{
					"--username",
					"admin",
					"--host",
					"https://love.aks1.eastus2.azure.cratedb.net:4200"
				}
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
