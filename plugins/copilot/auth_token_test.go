package copilot

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const exampleCopilotToken = "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE"

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: exampleCopilotToken,
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"COPILOT_GITHUB_TOKEN": exampleCopilotToken,
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"COPILOT_GITHUB_TOKEN": {
			Environment: map[string]string{
				"COPILOT_GITHUB_TOKEN": exampleCopilotToken,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: exampleCopilotToken,
					},
				},
			},
		},
		"GITHUB_TOKEN ignored": {
			Environment: map[string]string{
				"COPILOT_GITHUB_TOKEN": "",
				"GITHUB_TOKEN":         exampleCopilotToken,
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
	})
}
