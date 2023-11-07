package nomad

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "zwlwp719xssnxwvt3bq1rahafxsn0example",
				fieldname.Address: "http://localhost:4646",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NOMAD_TOKEN": "zwlwp719xssnxwvt3bq1rahafxsn0example",
					"NOMAD_ADDR": "http://localhost:4646",
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"NOMAD_TOKEN": "zwlwp719xssnxwvt3bq1rahafxsn0example",
				"NOMAD_ADDR": "http://nomad.example.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "zwlwp719xssnxwvt3bq1rahafxsn0example",
						fieldname.Address: "http://nomad.example.com",
					},
				},
			},
		},
	})
}
