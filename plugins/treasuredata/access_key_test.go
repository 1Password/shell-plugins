package treasuredata

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAccessKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "1/xxx",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"TD_API_KEY": "1/xxx",
				},
			},
		},
	})
}

func TestAccessKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessKey().Importer, map[string]plugintest.ImportCase{
		"TD config file": {
			Files: map[string]string{
				"~/.td/td.conf": plugintest.LoadFixture(t, "td.conf"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.User: "user@example.com",
						fieldname.APIKey: "1/xxx",
					},
				},
			},
		},
	})
}
