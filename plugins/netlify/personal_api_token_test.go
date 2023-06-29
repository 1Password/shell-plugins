package netlify

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAPIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "qaywqkexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NETLIFY_AUTH_TOKEN": "qaywqkexample",
				},
			},
		},
	})
}

func TestPersonalAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAPIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"NETLIFY_AUTH_TOKEN": "qaywqkexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "qaywqkexample",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in netlify/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "qaywqkexample",
				// 		},
				// 	},
			},
		},
	})
}
