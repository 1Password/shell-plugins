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
				fieldname.URL:    "https://api.prefect.cloud/api/accounts/be3b3f12-6de6-44ac-9228-1eefb857093c/workspaces/548cf0bf-f21d-456b-92a2-286233877ab9",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PREFECT_API_KEY": "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
					"PREFECT_API_URL": "https://api.prefect.cloud/api/accounts/be3b3f12-6de6-44ac-9228-1eefb857093c/workspaces/548cf0bf-f21d-456b-92a2-286233877ab9",
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
				"PREFECT_API_URL": "https://api.prefect.cloud/api/accounts/be3b3f12-6de6-44ac-9228-1eefb857093c/workspaces/548cf0bf-f21d-456b-92a2-286233877ab9",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
						fieldname.URL:    "https://api.prefect.cloud/api/accounts/be3b3f12-6de6-44ac-9228-1eefb857093c/workspaces/548cf0bf-f21d-456b-92a2-286233877ab9",
					},
				},
			},
		},
		"configfile in homedir": {
			Files: map[string]string{
				"~/.prefect/profiles.toml": plugintest.LoadFixture(t, "profiles.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:    "https://api.prefect.cloud/api/accounts/be3b3f12-6de6-44ac-9228-1eefb857093c/workspaces/548cf0bf-f21d-456b-92a2-286233877ab9",
						fieldname.APIKey: "pnu_dOXnQZOBq4Sst6hIesrdcbYVbnu1XEXAMPLE",
					},
				},
			},
		},
	})

}
