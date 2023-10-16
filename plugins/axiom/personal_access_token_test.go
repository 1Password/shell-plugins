package axiom

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "xapt-wovexreez0qf7zvkn935na41cudk2example",
				fieldname.Organization: "example",
				fieldname.Deployment:   "cloud",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AXIOM_TOKEN":      "xapt-wovexreez0qf7zvkn935na41cudk2example",
					"AXIOM_ORG_ID":     "example",
					"AXIOM_DEPLOYMENT": "cloud",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"AXIOM_TOKEN":      "xapt-wovexreez0qf7zvkn935na41cudk2example",
				"AXIOM_ORG_ID":     "example",
				"AXIOM_DEPLOYMENT": "cloud",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "xapt-wovexreez0qf7zvkn935na41cudk2example",
						fieldname.Organization: "example",
						fieldname.Deployment:   "cloud",
					},
				},
			},
		},
	})
}
