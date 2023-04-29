package st2

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestUserPassImporter(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.Username: "janedoe",
		fieldname.Password: "hunter2",
		fieldname.Website: "https://stackstorm.example.com",
	}
	plugintest.TestImporter(t, UserPass().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.st2/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields, NameHint: ""},
			},
		},
	})
}

func TestUserPassProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserPass().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "janedoe",
				fieldname.Password: "hunter2",
				fieldname.Website: "https://stackstorm.example.com",
			},
			CommandLine: []string{"st2"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"st2", "--config-file=~/.st2/config", "login", "janedoe"},
				Files: map[string]sdk.OutputFile{
					"~/.st2/config": {
						Contents: []byte(plugintest.LoadFixture(t, "config")),
					},
				},
			},
		},
	})
}
