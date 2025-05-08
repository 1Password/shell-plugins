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
				fieldname.APIKey: "2R6HTK2HeEpqPrQX4uGmJr736ng2MlBwA0tfiEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "2R6HTK2HeEpqPrQX4uGmJr736ng2MlBwA0tfiEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"IBMCLOUD_API_KEY": "2R6HTK2HeEpqPrQX4uGmJr736ng2MlBwA0tfiEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "2R6HTK2HeEpqPrQX4uGmJr736ng2MlBwA0tfiEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files:              map[string]string{},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
	})
}
