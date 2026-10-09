package qodana

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestProjectTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, ProjectToken().Importer, map[string]plugintest.ImportCase{
		"QODANA_TOKEN environment variable": {
			Environment: map[string]string{
				"QODANA_TOKEN": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.EXAMPLE",
					},
				},
			},
		},
	})
}

func TestProjectTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ProjectToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"QODANA_TOKEN": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.EXAMPLE",
				},
			},
		},
	})
}
