package spacelift

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Endpoint: "https://end.point",
				fieldname.APIKeyID: "abc123",
				fieldname.APIKeySecret: "def456",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SPACELIFT_API_KEY_ENDPOINT": "https://end.point",
					"SPACELIFT_API_KEY_ID": "abc123",
					"SPACELIFT_API_KEY_SECRET": "def456",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SPACELIFT_API_KEY_ENDPOINT": "xhdw",
				"SPACELIFT_API_KEY_ID": "abc123",
				"SPACELIFT_API_KEY_SECRET": "def456",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Endpoint: "xhdw",
						fieldname.APIKeyID: "abc123",
						fieldname.APIKeySecret: "def456",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.spacelift/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "dev",
					Fields: map[sdk.FieldName]string{
						fieldname.Endpoint: "https://mycorp.app.spacelift.io",
					},
				},
			},
		},
	})
}
