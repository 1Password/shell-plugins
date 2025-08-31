package anthropic

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
				fieldname.APIKey: "sk-ant-api03-6ZuKEYkQUs5tduCvwiQ7YmfDlx2h13q3ovHLfHX2NAXiUA81ixbluqFo4CfUqmodkho4HvYAmlYk0wuIFTLV5w-EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ANTHROPIC_API_KEY": "sk-ant-api03-6ZuKEYkQUs5tduCvwiQ7YmfDlx2h13q3ovHLfHX2NAXiUA81ixbluqFo4CfUqmodkho4HvYAmlYk0wuIFTLV5w-EXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"ANTHROPIC_API_KEY": "sk-ant-api03-6ZuKEYkQUs5tduCvwiQ7YmfDlx2h13q3ovHLfHX2NAXiUA81ixbluqFo4CfUqmodkho4HvYAmlYk0wuIFTLV5w-EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "sk-ant-api03-6ZuKEYkQUs5tduCvwiQ7YmfDlx2h13q3ovHLfHX2NAXiUA81ixbluqFo4CfUqmodkho4HvYAmlYk0wuIFTLV5w-EXAMPLE",
					},
				},
			},
		},
	})
}
