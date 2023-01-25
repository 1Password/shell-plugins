package ngrok

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.AuthToken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
				fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NGROK_AUTHTOKEN": "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
					"NGROK_API_KEY":   "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
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
						fieldname.AuthToken: "uSuQ7LUOJLs4xRbIySZ15F4v5KxfTnMknMdFEXAMPLE",
						fieldname.APIKey:    "L4STpMP3K8FNaQjBo5EAsXA2SThzq0J7BKD3jUZgtEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.ngrok2/ngrok.yml": plugintest.LoadFixture(t, "ngrok.yml"),
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
