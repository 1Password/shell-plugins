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
				fieldname.User_Access_Token: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
				fieldname.Endpoint: "https://huggingface.co",
				fieldname.API_URL: "https://api-inference.huggingface.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"HUGGINGFACE_TOKEN": "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
					"HF_ENDPOINT": "https://huggingface.co",
					"HF_INFERENCE_ENDPOINT": "https://api-inference.huggingface.com",

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
				"HF_ENDPOINT": "https://huggingface.co",
				"HF_INFERENCE_ENDPOINT": "https://api-inference.huggingface.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.UserAccessToken: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
						fieldname.Endpoint: "https://huggingface.co",
						fieldname.APIUrl: "https://api-inference.huggingface.com",
					},
				},
			},
		},
		"token file": {
			Files: map[string]string{
				"~/.cache/huggingface/token": plugintest.LoadFixture(t, "token"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "balaji_ceg@outlook.com",
					Fields: map[sdk.FieldName]string{
						fieldname.UserAccessToken: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
					},
				},
				{
					NameHint: "balaji_ceg@outlook.com",
					Fields: map[sdk.FieldName]string{
						fieldname.UserAccessToken: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
					},
				},
			},
		},
	})
}
