package huggingface

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
				fieldname.Token: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"HUGGINGFACE_TOKEN": "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ 
				"HUGGINGFACE_TOKEN": "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
					},
				},
			},
		},
	})
}
