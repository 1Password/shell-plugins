package okta

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OKTA_CLIENT_TOKEN":  "dIzt9kbedfNLtBNvWaprp39MaffIVjWxkZBEXAMPLE",
				"OKTA_CLIENT_ORGURL": "https://acme.okta.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:  "dIzt9kbedfNLtBNvWaprp39MaffIVjWxkZBEXAMPLE",
						fieldname.OrgURL: "https://acme.okta.com",
					},
				},
			},
		},
		"Okta config file": {
			Files: map[string]string{
				"~/.okta/okta.yaml": plugintest.LoadFixture(t, "okta.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:  "dIzt9kbedfNLtBNvWaprp39MaffIVjWxkZBEXAMPLE",
						fieldname.OrgURL: "https://acme.okta.com",
					},
				},
			},
		},
	})
}

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:  "dIzt9kbedfNLtBNvWaprp39MaffIVjWxkZBEXAMPLE",
				fieldname.OrgURL: "https://acme.okta.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OKTA_CLIENT_TOKEN":  "dIzt9kbedfNLtBNvWaprp39MaffIVjWxkZBEXAMPLE",
					"OKTA_CLIENT_ORGURL": "https://acme.okta.com",
				},
			},
		},
	})
}
