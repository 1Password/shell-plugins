package contentful

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "lDoY0qLbGt9Se8K71O42xI3y6XPxNSfFop4cjHHq1CEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CONTENTFUL_TOKEN": "lDoY0qLbGt9Se8K71O42xI3y6XPxNSfFop4cjHHq1CEXAMPLE",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.contentfulrc.json": plugintest.LoadFixture(t, "contentfulrc.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "CFPAT-XArJuIgiqjflHUGePPYbILIKOUzdFv2jTJMNEXAMPLE",
					},
				},
			},
		},
	})
}
