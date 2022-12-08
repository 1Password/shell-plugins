package snyk

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"SNYK_TOKEN": "etacgrrwj86t255ckguircl3kw3ftexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "etacgrrwj86t255ckguircl3kw3ftexample",
					},
				},
			},
		},
		"config file default path": {
			Files: map[string]string{
				"~/.config/configstore/snyk.json": plugintest.LoadFixture(t, "snyk.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "38j9ss3m5m3c44vi916z6p33m21xgexample",
					},
				},
			},
		},
	})
}

func TestAPITokenImporterFromCustomLocation(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "~/.snyk")

	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"config file custom path": {
			Files: map[string]string{
				"~/.snyk/configstore/snyk.json": plugintest.LoadFixture(t, "snyk.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "38j9ss3m5m3c44vi916z6p33m21xgexample",
					},
				},
			},
		},
	})
}

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "etacgrrwj86t255ckguircl3kw3ftexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNYK_TOKEN": "etacgrrwj86t255ckguircl3kw3ftexample",
				},
			},
		},
	})
}
