package prefect

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PREFECT_API_KEY": "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
				},
			},
		},
	})
}

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"PREFECT_API_KEY": "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in prefect/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
				// 		},
				// 	},
			},
		},
	})
}
