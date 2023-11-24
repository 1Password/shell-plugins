package exercism

import (
	"strings"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.URL:              "https://api.exercism.io/v1",
				fieldname.APIKey:           "v1o2p80wuf2qhnurrvf8rigro6sp38example",
				sdk.FieldName("Directory"): "/Users/username/exercism",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.config/exercism/user.json": {
						Contents: []byte(strings.Join(strings.Fields(plugintest.LoadFixture(t, "user.json")), "")),
					},
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.config/exercism/user.json": plugintest.LoadFixture(t, "user.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.URL:              "https://api.exercism.io/v1",
						fieldname.APIKey:           "v1o2p80wuf2qhnurrvf8rigro6sp38example",
						sdk.FieldName("Directory"): "/Users/username/exercism",
					},
				},
			},
		},
	})
}
