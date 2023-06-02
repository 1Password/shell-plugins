package clilol

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Address: "SyzkKWR6Pb4zspdKyji0ZFqSqw2v832uOUkaverJzsmwL5kqImIhD4XREXAMPLE",
				fieldname.Email:   "example@example.com",
				fieldname.APIKey:  "tnoe59h3w6jiv4tp0x9o72szmexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CLILOL_ADDRESS": "SyzkKWR6Pb4zspdKyji0ZFqSqw2v832uOUkaverJzsmwL5kqImIhD4XREXAMPLE",
					"CLILOL_EMAIL":   "example@example.com",
					"CLILOL_APIKEY":  "tnoe59h3w6jiv4tp0x9o72szmexample",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"CLILOL_ADDRESS": "SyzkKWR6Pb4zspdKyji0ZFqSqw2v832uOUkaverJzsmwL5kqImIhD4XREXAMPLE",
				"CLILOL_EMAIL":   "example@example.com",
				"CLILOL_APIKEY":  "tnoe59h3w6jiv4tp0x9o72szmexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: "SyzkKWR6Pb4zspdKyji0ZFqSqw2v832uOUkaverJzsmwL5kqImIhD4XREXAMPLE",
						fieldname.Email:   "example@example.com",
						fieldname.APIKey:  "tnoe59h3w6jiv4tp0x9o72szmexample",
					},
				},
			},
		},
	})
}
