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
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.APIKey: "rnd_Z7xMKp4NX1FoQNRyBpZs9yxDbu3i",
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"RENDER_API_KEY": "rnd_Z7xMKp4NX1FoQNRyBpZs9yxDbu3i",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "rnd_Z7xMKp4NX1FoQNRyBpZs9yxDbu3i",
					},
				},
			},
		},
		// TODO: If you implemented a config file importer, add a test file example in render/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				".render/config.yaml": plugintest.LoadFixture(t, "config.yml"),
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
