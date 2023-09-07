package vertica

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
				fieldname.Username: "vertica",
				fieldname.Password: "Il0v3y04",
				fieldname.Host:     "localhost",
				fieldname.Port:     "5433",
				fieldname.Database: "VMart",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"VSQL_USER":     "vertica",
					"VSQL_PASSWORD": "Il0v3y04",
					"VSQL_HOST":     "localhost",
					"VSQL_PORT":     "5433",
					"VSQL_DATABASE": "VMart",
				},
			},
		},
	})
}

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"VSQL_USER":     "vertica",
				"VSQL_PASSWORD": "Il0v3y04",
				"VSQL_HOST":     "localhost",
				"VSQL_PORT":     "5433",
				"VSQL_DATABASE": "VMart",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Username: "vertica",
						fieldname.Password: "Il0v3y04",
						fieldname.Host:     "localhost",
						fieldname.Port:     "5433",
						fieldname.Database: "VMart",
					},
				},
			},
		},
	})
}
