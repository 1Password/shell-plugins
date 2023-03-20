package ngrok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
			},
			CommandLine: []string{"ngrok"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"ngrok", "--config", "/tmp/config.yml"},
				Files: map[string]sdk.OutputFile{
					"/tmp/config.yml": {
						Contents: []byte(plugintest.LoadFixture(t, "config.yml")),
					},
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"NGROK_AUTHTOKEN": "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
					},
				},
			},
		},
		"config file for macos": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Application Support/ngrok/ngrok.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
					},
				},
			},
		},
		"config file for linux": {
			OS: "linux",
			Files: map[string]string{
				"~/.config/ngrok/ngrok.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
					},
				},
			},
		},
	})
}
