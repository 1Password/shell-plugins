package buildkite

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
				fieldname.Organization: "example",
				fieldname.Token: "bkua_abcdefghijklmnopqrstuvwxyz1234567890abcd",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BUILDKITE_ORGANIZATION_SLUG": "example",
					"BUILDKITE_API_TOKEN": "bkua_abcdefghijklmnopqrstuvwxyz1234567890abcd",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"BUILDKITE_ORGANIZATION_SLUG": "example",
				"BUILDKITE_API_TOKEN": "bkua_abcdefghijklmnopqrstuvwxyz1234567890abcd",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Organization: "example",
						fieldname.Token: "bkua_abcdefghijklmnopqrstuvwxyz1234567890abcd",
					},
				},
			},
		},
		"config file default path": {
			Files: map[string]string{
				"~/.config/bk.yaml": plugintest.LoadFixture(t, "bk.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Organization: "example",
						fieldname.Token: "bkua_abcdefghijklmnopqrstuvwxyz1234567890abcd",
					},
					NameHint: "example",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Organization: "example2",
						fieldname.Token: "bkua_abcdefghijklmnopqrstuvwxyz1234567890defg",
					},
					NameHint: "example2",
				},
			},
		},
	})
}

func TestAPIKeyNeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, BuildkiteCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"no for --help": {
			Args: []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"no for --version": {
			Args: []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"no for configure": {
			Args: []string{"configure"},
			ExpectedNeedsAuth: false,
		},
		"no for without args": {
			Args: []string{},
			ExpectedNeedsAuth: false,
		},
		"yes for all other commands": {
			Args: []string{"example"},
			ExpectedNeedsAuth: true,
		},
	})
}
