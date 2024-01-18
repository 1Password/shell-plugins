package sqitch

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "nBI8u8aF10TvQFfBlMedCDuEXAMPLE",
				fieldname.Username: "example_username",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SQITCH_PASSWORD": "nBI8u8aF10TvQFfBlMedCDuEXAMPLE",
					"SQITCH_USERNAME": "example_username",
				},
			},
		},
	})
}

func TestDatabaseCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SQITCH_PASSWORD": "nBI8u8aF10TvQFfBlMedCDuEXAMPLE",
				"SQITCH_USERNAME": "example_username",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Password: "nBI8u8aF10TvQFfBlMedCDuEXAMPLE",
						fieldname.Username: "example_username",
					},
				},
			},
		},
	})
}
