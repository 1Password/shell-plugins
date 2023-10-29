package shopify

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestCLITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, CLIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "atkn_0mgpyi6brmpxduriv7ukpp9r3lcjgsld7357svtw5fey5vlyriyauwxhgexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SHOPIFY_TOKEN": "atkn_0mgpyi6brmpxduriv7ukpp9r3lcjgsld7357svtw5fey5vlyriyauwxhgexample",
				},
			},
		},
	})
}

func TestCLITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, CLIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"SHOPIFY_TOKEN": "atkn_0mgpyi6brmpxduriv7ukpp9r3lcjgsld7357svtw5fey5vlyriyauwxhgexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "atkn_0mgpyi6brmpxduriv7ukpp9r3lcjgsld7357svtw5fey5vlyriyauwxhgexample",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in shopify/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "atkn_0mgpyi6brmpxduriv7ukpp9r3lcjgsld7357svtw5fey5vlyriyauwxhgexample",
			// 		},
			// 	},
			},
		},
	})
}
