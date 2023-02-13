package ibmcloud

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
				fieldname.APIKey: "tjLKgtZ5MSkC9zhSVGCeSfbxBbqr7KCbkfFaHEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "tjLKgtZ5MSkC9zhSVGCeSfbxBbqr7KCbkfFaHEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"IBMCLOUD_API_KEY": "tjLKgtZ5MSkC9zhSVGCeSfbxBbqr7KCbkfFaHEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "tjLKgtZ5MSkC9zhSVGCeSfbxBbqr7KCbkfFaHEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.bluemix/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "tjLKgtZ5MSkC9zhSVGCeSfbxBbqr7KCbkfFaHEXAMPLE",
					},
				},
			},
		},
	})
}
