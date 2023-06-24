package nirmata

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:   "fw90xpbsq8d1nzmdmbie0bcwk99x9rodqx72wfwif7y2hfbhq3gjg4bhcluw8b5qto5hwzfagsztgibbjs4rm8vswu6ppez8za9vhc2ozv5trexample",
				fieldname.Address: "https://nirmata.io",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NIRMATA_TOKEN": "fw90xpbsq8d1nzmdmbie0bcwk99x9rodqx72wfwif7y2hfbhq3gjg4bhcluw8b5qto5hwzfagsztgibbjs4rm8vswu6ppez8za9vhc2ozv5trexample",
					"NIRMATA_URL":   "https://nirmata.io",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"NIRMATA_TOKEN": "fw90xpbsq8d1nzmdmbie0bcwk99x9rodqx72wfwif7y2hfbhq3gjg4bhcluw8b5qto5hwzfagsztgibbjs4rm8vswu6ppez8za9vhc2ozv5trexample",
				"NIRMATA_URL":   "https://nirmata.io",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:   "fw90xpbsq8d1nzmdmbie0bcwk99x9rodqx72wfwif7y2hfbhq3gjg4bhcluw8b5qto5hwzfagsztgibbjs4rm8vswu6ppez8za9vhc2ozv5trexample",
						fieldname.Address: "https://nirmata.io",
					},
				},
			},
		},

		"config file": {
			Files: map[string]string{
				"~/.nirmata/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:   "fw90xpbsq8d1nzmdmbie0bcwk99x9rodqx72wfwif7y2hfbhq3gjg4bhcluw8b5qto5hwzfagsztgibbjs4rm8vswu6ppez8za9vhc2ozv5trexample",
						fieldname.Address: "https://nirmata.io",
						fieldname.Email:   "user@email.com",
					},
				},
			},
		},
	})
}
