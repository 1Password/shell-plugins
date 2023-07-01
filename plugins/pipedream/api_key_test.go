package pipedream

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"config file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key:   "9cfvvd7bp6099paodoua5shfoexample",
				fieldname.OrgID: "b_EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/pipedream/config": {
						Contents: []byte(plugintest.LoadFixture(t, "config")),
					},
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.config/pipedream/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "9cfvvd7bp6099paodoua5shfoexample",
						fieldname.OrgID:  "b_EXAMPLE",
					},
				},
			},
		},
	})
}
