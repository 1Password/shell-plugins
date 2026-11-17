package opencode

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const exampleAPIKey = "sk-0123456789abcdefghijEXAMPLE"

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: exampleAPIKey,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OPENCODE_API_KEY": exampleAPIKey,
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OPENCODE_API_KEY": exampleAPIKey,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: exampleAPIKey,
					},
				},
			},
		},
	})
}
