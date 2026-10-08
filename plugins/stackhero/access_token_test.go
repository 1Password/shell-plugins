package stackhero

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "usr-eqovniewliexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"STACKHERO_TOKEN": "usr-eqovniewliexample",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"STACKHERO_TOKEN": "usr-eqovniewliexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "usr-eqovniewliexample",
					},
				},
			},
		},
	})
}
