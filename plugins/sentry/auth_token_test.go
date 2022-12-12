package sentry

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"SENTRY_AUTH_TOKEN": "hmcxn4gjv9etm096p4v9ttxkvhj4tdm6ft6qmaj4szbb62bwu6mrl0gopexample",
				"SENTRY_ORG":        "acme",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "hmcxn4gjv9etm096p4v9ttxkvhj4tdm6ft6qmaj4szbb62bwu6mrl0gopexample",
						fieldname.Organization: "acme",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.sentryclirc": plugintest.LoadFixture(t, ".sentryclirc"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "mw4ms9tx4dci52bfr19sbj40lb0pu9w4camnf8w3hfzl8hckvkdocd28nexample",
						fieldname.Organization: "acme",
					},
				},
			},
		},
	})
}

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "hmcxn4gjv9etm096p4v9ttxkvhj4tdm6ft6qmaj4szbb62bwu6mrl0gopexample",
				fieldname.Organization: "acme",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SENTRY_AUTH_TOKEN": "hmcxn4gjv9etm096p4v9ttxkvhj4tdm6ft6qmaj4szbb62bwu6mrl0gopexample",
					"SENTRY_ORG":        "acme",
				},
			},
		},
	})
}
