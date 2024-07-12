package ngrok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ngrokProvisioner{}, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Authtoken: "cxG2Im21Yzkh8VnvFQaetlPHcQ9ZDUUk1IzzyHhcGcEXAMPLE",
				fieldname.APIKey:    "NQdxymVXmWC15916Mmy1vYkpzzNG6a84Bo4mYKuDahEXAMPLE",
			},
			CommandLine: []string{"ngrok"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"ngrok", "--config=/tmp/config.yml"},
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
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"NGROK_AUTHTOKEN": "cxG2Im21Yzkh8VnvFQaetlPHcQ9ZDUUk1IzzyHhcGcEXAMPLE",
				"NGROK_API_KEY":   "NQdxymVXmWC15916Mmy1vYkpzzNG6a84Bo4mYKuDahEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Authtoken: "cxG2Im21Yzkh8VnvFQaetlPHcQ9ZDUUk1IzzyHhcGcEXAMPLE",
						fieldname.APIKey:    "NQdxymVXmWC15916Mmy1vYkpzzNG6a84Bo4mYKuDahEXAMPLE",
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
						fieldname.Authtoken: "cxG2Im21Yzkh8VnvFQaetlPHcQ9ZDUUk1IzzyHhcGcEXAMPLE",
						fieldname.APIKey:    "NQdxymVXmWC15916Mmy1vYkpzzNG6a84Bo4mYKuDahEXAMPLE",
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
						fieldname.Authtoken: "cxG2Im21Yzkh8VnvFQaetlPHcQ9ZDUUk1IzzyHhcGcEXAMPLE",
						fieldname.APIKey:    "NQdxymVXmWC15916Mmy1vYkpzzNG6a84Bo4mYKuDahEXAMPLE",
					},
				},
			},
		},
	})
}

func TestGetAndReplaceConfigFlag(t *testing.T) {

	config, args := getConfigValueAndNewArgs([]string{}, "/newPath/to/newFile.json")
	assert.Equal(t, "", config)
	assert.Equal(t, []string{"--config=/newPath/to/newFile.json"}, args)

	config, args = getConfigValueAndNewArgs([]string{"--cache", "false", "--session", "asdefg345reger"}, "/newPath/to/newFile.json")
	assert.Equal(t, "", config)
	assert.Equal(t, []string{"--cache", "false", "--session", "asdefg345reger", "--config=/newPath/to/newFile.json"}, args)

	config, args = getConfigValueAndNewArgs([]string{"--cache", "false", "--config"}, "/newPath/to/newFile.json")
	assert.Equal(t, "", config)
	assert.Equal(t, []string{"--cache", "false", "--config", "/newPath/to/newFile.json"}, args)

	config, args = getConfigValueAndNewArgs([]string{"--cache", "false", "--config", "/path/to/file.json"}, "/newPath/to/newFile.json")
	assert.Equal(t, "/path/to/file.json", config)
	assert.Equal(t, []string{"--cache", "false", "--config", "/newPath/to/newFile.json"}, args)

	config, args = getConfigValueAndNewArgs([]string{"--config=/path/to/file.json", "--session", "asdefg345reger"}, "/newPath/to/newFile.json")
	assert.Equal(t, "/path/to/file.json", config)
	assert.Equal(t, []string{"--config=/newPath/to/newFile.json", "--session", "asdefg345reger"}, args)
}
