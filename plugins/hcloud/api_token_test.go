package hcloud

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestImportingTheHetznerConfig(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.Token: "dcAuOpQaCNvjzsNPmeGXvegHBdq4Zamx8QjI8ibxfErzy34fjL4ZOITFvdP5SKct",
	}

	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"Hcloud config file": {
			Files: map[string]string{
				"~/.config/hcloud/cli.toml": plugintest.LoadFixture(t, "hcloud.toml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields, NameHint: ""},
			},
		},
	})
}
