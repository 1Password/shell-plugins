package localxpose

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessToken: "PROVISIONERqLtcqQ8a3oRVfK5tiHzDOhEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"LX_ACCESS_TOKEN": "PROVISIONERqLtcqQ8a3oRVfK5tiHzDOhEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.AccessToken: "31QJpgl8FB9qLtcqQ8a3oRVfK5tiHzDOhEXAMPLE",
	}

	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"LX_ACCESS_TOKEN": "31QJpgl8FB9qLtcqQ8a3oRVfK5tiHzDOhEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: expectedFields,
				},
			},
		},
		"file": {
			Files: map[string]string{
				"~/.localxpose/.access": plugintest.LoadFixture(t, "localxpose.access"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: expectedFields,
				},
			},
		},
	})
}
