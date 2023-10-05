package docker

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "david",
				fieldname.Secret:   "passw0rd1",
				fieldname.Host:     "https://index.docker.io/v1",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{
					"https://index.docker.io/v1", "--username", "david", "--password", "passw0rd1",
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.docker/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "david",
						fieldname.Secret:   "passw0rd1",
						fieldname.Host:     "https://index.docker.io/v1",
					},
				},
			},
		},
		"no url config file": {
			Files: map[string]string{
				"~/.docker/config.json": plugintest.LoadFixture(t, "no_url_config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "user123",
						fieldname.Secret:   "passw0rd2",
						fieldname.Host:     "",
					},
				},
			},
		},
	})
}
