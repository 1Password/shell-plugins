package openbao

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Token().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string {
				fieldname.Token: "s.UrpjvNwnaPjTFFj2RAyEXAMPLE",
				fieldname.Address: "https://bao.acme.com",
				fieldname.Namespace: "default",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"BAO_TOKEN":     "s.UrpjvNwnaPjTFFj2RAyEXAMPLE",
					"BAO_ADDR":      "https://bao.acme.com",
					"BAO_NAMESPACE": "default",
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, Token().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string {
				"BAO_TOKEN":     "s.UrpjvNwnaPjTFFj2RAyEXAMPLE",
				"BAO_ADDR":      "https://bao.acme.com",
				"BAO_NAMESPACE": "default",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "s.UrpjvNwnaPjTFFj2RAyEXAMPLE",
						fieldname.Address: "https://bao.acme.com",
						fieldname.Namespace: "default",
					},
				},
			},
		},
	})
}
