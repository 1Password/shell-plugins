package docker

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestUserCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.: "7nJeRV69EYA3zvEjpSmwb5GEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DOCKER": "7nJeRV69EYA3zvEjpSmwb5GEXAMPLE",
				},
			},
		},
	})
}

func TestUserCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, UserCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"DOCKER": "7nJeRV69EYA3zvEjpSmwb5GEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.: "7nJeRV69EYA3zvEjpSmwb5GEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in docker/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "7nJeRV69EYA3zvEjpSmwb5GEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
