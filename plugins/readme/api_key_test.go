package readme

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"RDME_API_KEY": "rdme_9o50rxz28p0msun40apgyzvkdji2vhxyd7b8emioclkrx57ucpb5x2d31yu39taexample",
				"RDME_EMAIL":   "owlbert@readme.io",
				"RDME_PROJECT": "test-subdomain",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "owlbert@readme.io",
						fieldname.APIKey:   "rdme_9o50rxz28p0msun40apgyzvkdji2vhxyd7b8emioclkrx57ucpb5x2d31yu39taexample",
						fieldname.Website:  "test-subdomain",
					},
				},
			},
		},
		"ReadMe config file": {
			Files: map[string]string{
				"~/.config/configstore/rdme-production.json": plugintest.LoadFixture(t, "readme.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "owlbert@readme.io",
						fieldname.APIKey:   "rdme_9o50rxz28p0msun40apgyzvkdji2vhxyd7b8emioclkrx57ucpb5x2d31yu39taexample",
						fieldname.Website:  "https://dash.readme.com/go/test-subdomain",
					},
					NameHint: "test-subdomain",
				},
			},
		},
	})
}

func TestSecretKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "owlbert@readme.io",
				fieldname.APIKey:   "rdme_9o50rxz28p0msun40apgyzvkdji2vhxyd7b8emioclkrx57ucpb5x2d31yu39taexample",
				fieldname.Website:  "https://dash.readme.com/go/test-subdomain",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"RDME_EMAIL":   "owlbert@readme.io",
					"RDME_API_KEY": "rdme_9o50rxz28p0msun40apgyzvkdji2vhxyd7b8emioclkrx57ucpb5x2d31yu39taexample",
					"RDME_PROJECT": "https://dash.readme.com/go/test-subdomain",
				},
			},
		},
	})
}
