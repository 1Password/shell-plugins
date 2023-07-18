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
				fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
				fieldname.OrgID:  "YbEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/pipedream/config": {
						Contents: []byte(plugintest.LoadFixture(t, "provision")),
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
				"~/.config/pipedream/config": plugintest.LoadFixture(t, "import"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "ugvfxesz62ycsl42z49c0t1hjexample",
						fieldname.OrgID:  "YbEXAMPLE",
					},
					NameHint: "DEFAULT",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "5puf32rvhkz83c6oj4wpxvaniexample",
						fieldname.OrgID:  "KVEXAMPLE",
					},
					NameHint: "first",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "lgx1amb0qf7mjy6y7nkgfc3x9example",
					},
					NameHint: "second",
				},
			},
		},
	})
}
