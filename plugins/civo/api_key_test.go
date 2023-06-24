package civo

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
				fieldname.APIKey: "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
				fieldname.APIKeyID: "civoapikey",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CIVO_TOKEN": "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
					"CIVO_API_KEY": "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
					"CIVO_API_KEY_NAME": "civoapikey",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ 
				"CIVO_TOKEN": "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
				"CIVO_API_KEY": "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
				"CIVO_API_KEY_NAME": "civoapikey",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
						fieldname.APIKeyID: "civoapikey",
					},
				},
			},
		},
		
		"config file": {
			Files: map[string]string{
				"~/.civo.json": plugintest.LoadFixture(t, "civo.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "vvUmolh6TG7VMzCkhfWvOCJfpZMzM6W5hWkjurT63FoEXAMPLE",
						fieldname.APIKeyID: "civoapikey",
					},
				},
			},
		},
	})
}
