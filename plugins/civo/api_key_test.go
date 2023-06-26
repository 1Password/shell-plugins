package civo

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
				fieldname.APIKey:   "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
				fieldname.APIKeyID: "testdemoname",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CIVO_TOKEN":        "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
					"CIVO_API_KEY":      "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
					"CIVO_API_KEY_NAME": "testdemoname",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"CIVO_TOKEN":        "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
				"CIVO_API_KEY":      "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
				"CIVO_API_KEY_NAME": "testdemoname",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:   "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
						fieldname.APIKeyID: "testdemoname",
					},
				},
			},
		},

		"config file": {
			Files: map[string]string{

				"~/.civo.json": plugintest.LoadFixture(t, ".civo.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{

				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey:        "XFIx85McyfCQc490j1tBa5b5s2XiWerNdOdfnkrOnchEXAMPLE",
						fieldname.APIKeyID:      "testdemoname",
						fieldname.DefaultRegion: "LON1",
					},
				},
			},
		},
	})
}
