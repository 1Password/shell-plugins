package mongodbatlas

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.PublicKey:  "eexample",
				fieldname.PrivateKey: "qohcbhiu-26ag-wpwf-maqn-cw5xlexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"MONGODB_ATLAS_PUBLIC_API_KEY":  "eexample",
					"MONGODB_ATLAS_PRIVATE_API_KEY": "qohcbhiu-26ag-wpwf-maqn-cw5xlexample",
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"MONGODB_ATLAS_PUBLIC_API_KEY":  "eexample",
				"MONGODB_ATLAS_PRIVATE_API_KEY": "qohcbhiu-26ag-wpwf-maqn-cw5xlexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.PublicKey:  "eexample",
						fieldname.PrivateKey: "qohcbhiu-26ag-wpwf-maqn-cw5xlexample",
					},
				},
			},
		},
	})
}
