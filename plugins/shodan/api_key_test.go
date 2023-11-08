package shodan

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
				fieldname.APIKey: "ddXfzwQOIjTaaxGMzcxXYR6Q0EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/shodan/api_key": {Contents: []byte("ddXfzwQOIjTaaxGMzcxXYR6Q0EXAMPLE")},
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.config/shodan/api_key": plugintest.LoadFixture(t, "api_key"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "ddXfzwQOIjTaaxGMzcxXYR6Q0EXAMPLE",
					},
				},
			},
		},
	})
}
