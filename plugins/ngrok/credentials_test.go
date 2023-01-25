package ngrok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

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
						fieldname.AuthToken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
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
						fieldname.AuthToken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
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
						fieldname.AuthToken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
						fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
					},
				},
			},
		},
	})
}
