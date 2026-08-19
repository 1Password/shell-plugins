package netlify

import (
	"testing"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "tGtp-IMFGyRcoLdK40zQ4ENKfvDeIOASs1ilEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NETLIFY_TOKEN": "tGtp-IMFGyRcoLdK40zQ4ENKfvDeIOASs1ilEXAMPLE",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"NETLIFY_TOKEN": "tGtp-IMFGyRcoLdK40zQ4ENKfvDeIOASs1ilEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "tGtp-IMFGyRcoLdK40zQ4ENKfvDeIOASs1ilEXAMPLE",
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
			// 			fieldname.Token: "tGtp-IMFGyRcoLdK40zQ4ENKfvDeIOASs1ilEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}