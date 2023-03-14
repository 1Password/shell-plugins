package zendesk

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
				fieldname.OrgURL:   "subdomain",
				fieldname.Username: "wendy@appleseed.com",
				fieldname.Token:    "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ZENDESK_SUBDOMAIN": "subdomain",
					"ZENDESK_EMAIL":     "wendy@appleseed.com",
					"ZENDESK_API_TOKEN": "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"ZENDESK_SUBDOMAIN": "subdomain",
				"ZENDESK_EMAIL":     "wendy@appleseed.com",
				"ZENDESK_API_TOKEN": "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.OrgURL:   "subdomain",
						fieldname.Username: "wendy@appleseed.com",
						fieldname.Token:    "TPPmg1SEWr4fDGQhaUHsxETCUrBEIJKm0EXAMPLE",
					},
				},
			},
		},
	})
}
