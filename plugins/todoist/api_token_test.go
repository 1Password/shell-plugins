package todoist

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"config file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "dbq9y65uguqrk4ognfhdiwcc0zx34z20pexample",
			},
			CommandLine: []string{"todoist"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"todoist"},
				Files: map[string]sdk.OutputFile{
					"~/.config/todoist/config.json": {
						Contents: []byte(plugintest.LoadFixture(t, "config.json")),
					},
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.config/todoist/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dbq9y65uguqrk4ognfhdiwcc0zx34z20pexample",
					},
				},
			},
		},
	})
}
