package ngrok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
				fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
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
				"NGROK_AUTHTOKEN": "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
				"NGROK_API_KEY":   "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Authtoken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
						fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
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
						fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
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
						fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
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
