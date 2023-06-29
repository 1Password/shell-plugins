package huggingface

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(
		t, UserAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
			"default": {
				ItemFields: map[sdk.FieldName]string{
					fieldname.UserAccessToken: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
					fieldname.Endpoint:        "https://huggingface.co",
					fieldname.APIUrl:          "https://api-inference.huggingface.com",
				},
				ExpectedOutput: sdk.ProvisionOutput{
					Environment: map[string]string{
						"HUGGING_FACE_HUB_TOKEN": "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
						"HF_ENDPOINT":            "https://huggingface.co",
						"HF_INFERENCE_ENDPOINT":  "https://api-inference.huggingface.com",
					},
				},
			},
		})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(
		t, UserAccessToken().Importer, map[string]plugintest.ImportCase{
			"config file (macOS)": {
				Files: map[string]string{
					"~/.cache/huggingface/token": plugintest.LoadFixture(t, "token"),
				},
				ExpectedCandidates: []sdk.ImportCandidate{
					{
						Fields: map[sdk.FieldName]string{
							fieldname.UserAccessToken: "hf_yVvZeburdKtnwkVCWPXimmNwaFuEXAMPLE",
						},
					},
				},
			},
		},
	)
}
