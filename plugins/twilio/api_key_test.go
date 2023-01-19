package twilio

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
				"TWILIO_ACCOUNT_SID": "AC0XRE2037R1E7EHWTFD3PXAMSPEXAMPLE",
				"TWILIO_API_KEY":     "SK5CPBH8WCGB8Q1XFCBNQBKGYCZEXAMPLE",
				"TWILIO_API_SECRET":  "WKLK57W483X5YE74UN49WV8MJEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccountSID: "AC0XRE2037R1E7EHWTFD3PXAMSPEXAMPLE",
						fieldname.APIKey:     "SK5CPBH8WCGB8Q1XFCBNQBKGYCZEXAMPLE",
						fieldname.APISecret:  "WKLK57W483X5YE74UN49WV8MJEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.twilio-cli/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "prod",
					Fields: map[sdk.FieldName]string{
						fieldname.AccountSID: "ACtdgcL157CFBnWjM7ZeJKysIjEEXAMPLE",
						fieldname.APIKey:     "SK4bqP76ByZgGEuwqm0eTFzYrWBEXAMPLE",
						fieldname.APISecret:  "1KAe9Vshg4EkUvaBVS8pZwDS1EXAMPLE",
					},
				},
				{
					NameHint: "dev",
					Fields: map[sdk.FieldName]string{
						fieldname.AccountSID: "ACZBgJOUfaX2AuuLMWK7jT3tdS9EXAMPLE",
						fieldname.APIKey:     "SKrhNjOV2LgR1xatOpVFsMa5fOpEXAMPLE",
						fieldname.APISecret:  "4ELE8BqwbCrzbyTqu7HNylK00EXAMPLE",
					},
				},
			},
		},
	})
}

func TestSecretKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccountSID: "AC0XRE2037R1E7EHWTFD3PXAMSPEXAMPLE",
				fieldname.APIKey:     "SK5CPBH8WCGB8Q1XFCBNQBKGYCZEXAMPLE",
				fieldname.APISecret:  "WKLK57W483X5YE74UN49WV8MJEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"TWILIO_ACCOUNT_SID": "AC0XRE2037R1E7EHWTFD3PXAMSPEXAMPLE",
					"TWILIO_API_KEY":     "SK5CPBH8WCGB8Q1XFCBNQBKGYCZEXAMPLE",
					"TWILIO_API_SECRET":  "WKLK57W483X5YE74UN49WV8MJEXAMPLE",
				},
			},
		},
	})
}
