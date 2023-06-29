package firebase

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "dummy_firebase_3bhfuelt31a99503j251bua8rov58m2example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"FIREBASE_TOKEN": "dummy_firebase_3bhfuelt31a99503j251bua8rov58m2example",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"FIREBASE_TOKEN": "dummy_firebase_3bhfuelt31a99503j251bua8rov58m2example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dummy_firebase_3bhfuelt31a99503j251bua8rov58m2example",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in firebase/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "dummy_firebase_3bhfuelt31a99503j251bua8rov58m2example",
			// 		},
			// 	},
			},
		},
	})
}
