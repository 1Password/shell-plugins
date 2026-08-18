package grok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "xai-EXAMPLE1234abcdefEXAMPLE1234abcdef",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"XAI_API_KEY": "xai-EXAMPLE1234abcdefEXAMPLE1234abcdef",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"XAI_API_KEY": "xai-EXAMPLE1234abcdefEXAMPLE1234abcdef",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "xai-EXAMPLE1234abcdefEXAMPLE1234abcdef",
					},
				},
			},
		},
	})
}
