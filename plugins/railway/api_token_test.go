package railway

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "abcdefghijklm-1234567890-example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"RAILWAY_API_TOKEN": "abcdefghijklm-1234567890-example",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"RAILWAY_API_TOKEN": "abcdefghijklm-1234567890-example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "abcdefghijklm-1234567890-example",
					},
				},
			},
		},
	})
}
