package databricks

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const example_host = "https://myinstance.databricks.com"
const example_token = "dapif13ac4b49d1cb31f69f678e39602e381-2"

func TestDatabricksPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host:  example_host,
				fieldname.Token: example_token,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DATABRICKS_HOST":  example_host,
					"DATABRICKS_TOKEN": example_token,
				},
			},
		},
	})
}

func TestDatabricksPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"DATABRICKS_HOST":  example_host,
				"DATABRICKS_TOKEN": example_token,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host:  example_host,
						fieldname.Token: example_token,
					},
				},
			},
		},
	})

	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.databrickscfg": plugintest.LoadFixture(t, "databrickscfg"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "secondprofile",
					Fields: map[sdk.FieldName]string{
						fieldname.Host:  example_host,
						fieldname.Token: example_token,
					},
				},
				{
					NameHint: "thirdprofile",
					Fields: map[sdk.FieldName]string{
						fieldname.Host:  example_host,
						fieldname.Token: example_token,
					},
				},
			},
		},
	})
}
