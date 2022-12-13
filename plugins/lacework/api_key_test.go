package lacework

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"LW_ACCOUNT":    "example",
				"LW_API_KEY":    "EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
				"LW_API_SECRET": "_89368245c62f8d6d35e7c6626example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:   "example",
						fieldname.APIKeyID:  "EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
						fieldname.APISecret: "_89368245c62f8d6d35e7c6626example",
					},
				},
			},
		},
		"Lacework config file": {
			Files: map[string]string{
				"~/.lacework.toml": plugintest.LoadFixture(t, "lacework.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:   "example",
						fieldname.APIKeyID:  "EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
						fieldname.APISecret: "_89368245c62f8d6d35e7c6626example",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:   "example2",
						fieldname.APIKeyID:  "EXAMPLE2_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
						fieldname.APISecret: "_69ee92b3c71a6b27436a648acexample",
					},
					NameHint: "example",
				},
			},
		},
	})
}

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Account:   "example",
				fieldname.APIKeyID:  "EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
				fieldname.APISecret: "_89368245c62f8d6d35e7c6626example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"LW_ACCOUNT":    "example",
					"LW_API_KEY":    "EXAMPLE_1234567890ABCDE1EXAMPLE1EXAMPLE123456789EXAMPLE",
					"LW_API_SECRET": "_89368245c62f8d6d35e7c6626example",
				},
			},
		},
	})
}
