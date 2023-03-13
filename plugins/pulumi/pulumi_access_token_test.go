package pulumi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PulumiAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PULUMI_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PulumiAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"PULUMI_TOKEN": "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in pulumi/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				// 	{
				// 		Fields: map[sdk.FieldName]string{
				// 			fieldname.Token: "pul-8s9b3qf8rx7x8x8pn03ibkemilm1zfs10example",
				// 		},
				// 	},
			},
		},
	})
}
