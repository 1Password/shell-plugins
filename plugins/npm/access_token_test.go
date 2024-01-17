package npm

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
				fieldname.Token: "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NPM_TOKEN": "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"NPM_TOKEN": "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in npm/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "npm_4F1h0xp0lBn0XvZ4RGfpnpoawECBMEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
