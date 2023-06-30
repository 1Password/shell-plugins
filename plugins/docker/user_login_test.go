package docker

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestUserLoginProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserLogin().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "cheithanya",
				fieldname.Password: "Dockerforch@123",
				//	fieldname.URL:      "https://index.docker.io/v1/",
			},
			ExpectedOutput: sdk.ProvisionOutput{

				Files: map[string]sdk.OutputFile{
					"~/.docker/config.json": {Contents: []byte(plugintest.LoadFixture(t, "config.json"))},
				},
			},
		},
	})
}

func TestUserLoginImporter(t *testing.T) {
	plugintest.TestImporter(t, UserLogin().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.docker/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "cheithanya",
						fieldname.Password: "Dockerforch@123",
					},
				},
			},
		},
	})
}
