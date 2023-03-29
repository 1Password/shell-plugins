package awscdk

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Key: "YV4DFI6DLBV9KEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AWS_CDK_KEY": "YV4DFI6DLBV9KEXAMPLE",
				},
			},
		},
	})
}

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"AWS_CDK_KEY": "YV4DFI6DLBV9KEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "YV4DFI6DLBV9KEXAMPLE",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in awscdk/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "YV4DFI6DLBV9KEXAMPLE",
				// 		},
				// 	},
			},
		},
	})
}
