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
				fieldname.APIKey: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
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
						fieldname.APIKey: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
					},
				},
			},
		},
	})
}
