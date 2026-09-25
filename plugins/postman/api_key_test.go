package postman

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
				fieldname.APIKey: "jNMdtxuUHqiSMNmwYu3OXyhOgYKse6H82uhghe0Zw3K92ZEfXhL8wvLzX",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"POSTMAN_API_KEY": "jNMdtxuUHqiSMNmwYu3OXyhOgYKse6H82uhghe0Zw3K92ZEfXhL8wvLzX",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"POSTMAN_API_KEY": "jNMdtxuUHqiSMNmwYu3OXyhOgYKse6H82uhghe0Zw3K92ZEfXhL8wvLzX",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "jNMdtxuUHqiSMNmwYu3OXyhOgYKse6H82uhghe0Zw3K92ZEfXhL8wvLzX",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.postman/postmanrc": plugintest.LoadFixture(t, "postmanrc"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "jNMdtxuUHqiSMNmwYu3OXyhOgYKse6H82uhghe0Zw3K92ZEfXhL8wvLzX",
					},
				},
			},
		},
	})
}
