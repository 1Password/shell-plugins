package postgresql

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDatabaseCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"PGHOST":     "localhost",
				"PGPORT":     "5432",
				"PGUSER":     "root",
				"PGPASSWORD": "123456",
				"PGDATABASE": "test",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host:     "localhost",
						fieldname.Port:     "5432",
						fieldname.User:     "root",
						fieldname.Password: "123456",
						fieldname.Database: "test",
					},
				},
			},
		},
	})
}

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host:     "localhost",
				fieldname.Port:     "5432",
				fieldname.User:     "root",
				fieldname.Password: "123456",
				fieldname.Database: "test",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PGHOST":     "localhost",
					"PGPORT":     "5432",
					"PGUSER":     "root",
					"PGPASSWORD": "123456",
					"PGDATABASE": "test",
				},
			},
		},
	})
}
