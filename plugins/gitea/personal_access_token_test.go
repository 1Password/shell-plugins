package gitea

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
				fieldname.HostAddress: "https://git.example.com",
				fieldname.User:        "example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/tea/config.yml": {
						Contents: []byte(plugintest.LoadFixture(t, "config.yml")),
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{

		"config file": {
			Files: map[string]string{
				"~/.config/tea/config.yml": plugintest.LoadFixture(t, "import_config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
						fieldname.HostAddress: "https://git.example.com",
						fieldname.User:        "example",
					},
					NameHint: "git.example.com",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "enjkarzu2ca5ffcnvzaczxncuczeoq9utlpqqrzs",
						fieldname.HostAddress: "https://gitea.com",
						fieldname.User:        "example@example.com",
					},
					NameHint: "gitea.com",
				},
			},
		},
	})
}
