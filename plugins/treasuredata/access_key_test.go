package treasuredata

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "1/xxx",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"TD_API_KEY": "1/xxx",
				},
			},
		},
	})
}

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"env var TD_API_KEY": {
			Environment: map[string]string{
				"TD_API_KEY": "1/xxx",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "1/xxx",
					},
				},
			},
		},
		"env var TREASURE_DATA_API_KEY": {
			Environment: map[string]string{
				"TREASURE_DATA_API_KEY": "1/xxx",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "1/xxx",
					},
				},
			},
		},
		"TD config file": {
			Files: map[string]string{
				"~/.td/td.conf": plugintest.LoadFixture(t, "td.conf"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "1/xxx",
					},
				},
			},
		},
	})
}

func TestAPIKeyNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, TreasureDataCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no for -c": {
			Args:              []string{"-c"},
			ExpectedNeedsAuth: false,
		},
		"no for -k": {
			Args:              []string{"-k"},
			ExpectedNeedsAuth: false,
		},
	})
}
