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
				fieldname.APIKey: "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
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
