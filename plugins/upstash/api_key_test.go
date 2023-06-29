package upstash

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
				fieldname.APIKey: "d68850db-69f7-qxe9pubcmjnqfgyexample",
				fieldname.Email :  "fakememail12@gmail.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"UPSTASH_API_KEY": "d68850db-69f7-qxe9pubcmjnqfgyexample",
					"UPSTASH_EMAIL" : "fakememail12@gmail.com",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ 
				"UPSTASH_API_KEY": "d68850db-69f7-qxe9pubcmjnqfgyexample",
				"UPSTASH_EMAIL" : "fakememail12@gmail.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "d68850db-69f7-qxe9pubcmjnqfgyexample",
						fieldname.Email :  "fakememail12@gmail.com",
					},
				},
			},
		},
		
		"config file": {
			Files: map[string]string{
				"~/.upstash.json": plugintest.LoadFixture(t, ".upstash.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "d68850db-69f7-qxe9pubcmjnqfgyexample",
						fieldname.Email :  "fakememail12@gmail.com",
					},
				},
			},
		},
	})
}
