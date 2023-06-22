package appwrite

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCLITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, CLIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "dPmVEzGfQFzWGqakEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"APPWRITE_PASSWORD": "dPmVEzGfQFzWGqakEXAMPLE",
				},
			},
		},
	})
}

func TestCLITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, CLIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"APPWRITE_PASSWORD": "dPmVEzGfQFzWGqakEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "dPmVEzGfQFzWGqakEXAMPLE",
					},
				},
			},
		},
	})
}
