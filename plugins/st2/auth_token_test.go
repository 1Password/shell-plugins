package st2

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "i03ek4dcx61v1utonqbcnl1cjexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ST2_AUTH_TOKEN": "i03ek4dcx61v1utonqbcnl1cjexample",
				},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"ST2_AUTH_TOKEN": "i03ek4dcx61v1utonqbcnl1cjexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "i03ek4dcx61v1utonqbcnl1cjexample",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.st2/token": plugintest.LoadFixture(t, "token"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "i03ek4dcx61v1utonqbcnl1cjexample",
					},
				},
			},
		},
	})
}
