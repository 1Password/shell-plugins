package azure

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestServicePrincipalProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ServicePrincipal().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.ClientID:     "6a768b1b-8eee-40e9-81b9-c0bca485f33e",
				fieldname.ClientSecret: "T.L8Q~IhftYnbzhuCADzkfIRZ9p8e5kiyjSR4cjb",
				fieldname.TenantID:     "83f25d69-e104-45ba-8939-4d9103b03b3c",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"AZURE_CLIENT_ID":     "6a768b1b-8eee-40e9-81b9-c0bca485f33e",
					"AZURE_CLIENT_SECRET": "T.L8Q~IhftYnbzhuCADzkfIRZ9p8e5kiyjSR4cjb",
					"AZURE_TENANT_ID":     "83f25d69-e104-45ba-8939-4d9103b03b3c",
				},
			},
		},
	})
}

func TestServicePrincipalImporter(t *testing.T) {
	plugintest.TestImporter(t, ServicePrincipal().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"AZURE_CLIENT_ID":     "f3593130-88fe-4f93-95c5-1c62d2d56a36",
				"AZURE_CLIENT_SECRET": "51U8Q~tqPdbkcxgYChTSie9br1_gIgHh2ZeWKcHX",
				"AZURE_TENANT_ID":     "a16c4485-bc84-41a3-90f4-b00498389c06",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.ClientID:     "f3593130-88fe-4f93-95c5-1c62d2d56a36",
						fieldname.ClientSecret: "51U8Q~tqPdbkcxgYChTSie9br1_gIgHh2ZeWKcHX",
						fieldname.TenantID:     "a16c4485-bc84-41a3-90f4-b00498389c06",
					},
				},
			},
		},
	})
}
