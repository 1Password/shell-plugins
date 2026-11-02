package notion

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "ntn_47912285616b5Lg0KVxJ9ftlZsOd3QfWRB7mz6kEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NOTION_API_TOKEN": "ntn_47912285616b5Lg0KVxJ9ftlZsOd3QfWRB7mz6kEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"NOTION_API_TOKEN": "ntn_47912285616b5Lg0KVxJ9ftlZsOd3QfWRB7mz6kEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "ntn_47912285616b5Lg0KVxJ9ftlZsOd3QfWRB7mz6kEXAMPLE",
					},
				},
			},
		},
	})
}
