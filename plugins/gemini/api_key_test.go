package gemini

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
				fieldname.APIKey: "1y1RvqUVfm6eTlAUfzFfxcEo1EP9dWMdEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GEMINI_API_KEY": "1y1RvqUVfm6eTlAUfzFfxcEo1EP9dWMdEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"GEMINI_API_KEY": "1y1RvqUVfm6eTlAUfzFfxcEo1EP9dWMdEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "1y1RvqUVfm6eTlAUfzFfxcEo1EP9dWMdEXAMPLE",
					},
				},
			},
		},
	})
}
