package ibmcloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCLIProvisioner(t *testing.T) {
	provisioner := CLIProvisioner{}

	testCases := map[string]plugintest.ProvisionCase{
		"with api key only": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
			},
			CommandLine: []string{"ibmcloud", "login"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				},
				CommandLine: []string{"ibmcloud", "login"},
			},
		},
		"with api key and resource group": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey:                "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				sdk.FieldName("resource group"): "my-resource-group",
			},
			CommandLine: []string{"ibmcloud", "login"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				},
				CommandLine: []string{"ibmcloud", "login", "-g", "my-resource-group"},
			},
		},
		"with resource group but not login command": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey:                "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				sdk.FieldName("resource group"): "my-resource-group",
			},
			CommandLine: []string{"target"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				},
				CommandLine: []string{"target"},
			},
		},
		"with resource group and existing group flag": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey:                "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				sdk.FieldName("resource group"): "my-resource-group",
			},
			CommandLine: []string{"login", "-g", "existing-group"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"IBMCLOUD_API_KEY": "GeYS3RmGXo7cQhY8UboUSLmWarFF1HGqv4fVKEXAMPLE",
				},
				CommandLine: []string{"login", "-g", "existing-group"},
			},
		},
	}

	plugintest.TestProvisioner(t, provisioner, testCases)
}
