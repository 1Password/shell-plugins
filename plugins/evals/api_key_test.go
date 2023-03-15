package evals

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
				fieldname.APIKey: "sk-ZJipjm9euld2QcpLIX3NCDUiLgexQkvdn0WGlKayAEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OPENAI_API_KEY": "sk-ZJipjm9euld2QcpLIX3NCDUiLgexQkvdn0WGlKayAEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OPENAI_API_KEY": "sk-ZJipjm9euld2QcpLIX3NCDUiLgexQkvdn0WGlKayAEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "sk-ZJipjm9euld2QcpLIX3NCDUiLgexQkvdn0WGlKayAEXAMPLE",
					},
				},
			},
		},
	})
}
