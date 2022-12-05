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
