package render

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
				fieldname.APIKey: "rnd_Z7xMKp4NX1FoQNRyBpZs9yxDbu3i",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.render/config.yaml": {
						Contents: []byte(plugintest.LoadFixture(t, "config.yaml")),
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
				"~/.render/config.yaml": plugintest.LoadFixture(t, "config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "rnd_Z7xMKp4NX1FoQNRyBpZs9yxDbu3i",
					},
				},
			},
		},
	})
}
