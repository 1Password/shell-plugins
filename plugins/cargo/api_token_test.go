package cargo

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
				fieldname.Token: "9xAQsMIO2UubpsgD2eUOKqXEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CARGO_REGISTRY_TOKEN": "9xAQsMIO2UubpsgD2eUOKqXEXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"CARGO_REGISTRY_TOKEN": "9xAQsMIO2UubpsgD2eUOKqXEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "9xAQsMIO2UubpsgD2eUOKqXEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.cargo/credentials.toml": plugintest.LoadFixture(t, "credentials.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "9xAQsMIO2UubpsgD2eUOKqXEXAMPLE",
					},
				},
			},
		},
	})
}
