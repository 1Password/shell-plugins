package dbtredshift

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
				fieldname.Host:     "examplecluster.abc123xyz789.us-west-1.redshift.amazonaws.com",
				fieldname.Port:     "5439",
				fieldname.User:     "awsuser",
				fieldname.Password: "my_password",
				fieldname.Database: "dev",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DBT_HOST":     "examplecluster.abc123xyz789.us-west-1.redshift.amazonaws.com",
					"DBT_PORT":     "5439",
					"DBT_USER":     "awsuser",
					"DBT_PASSWORD": "my_password",
					"DBT_DB":       "dev",
				},
			},
		},
	})
}

func TestDatabaseCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"DBT_HOST":     "examplecluster.abc123xyz789.us-west-1.redshift.amazonaws.com",
				"DBT_PORT":     "5439",
				"DBT_USER":     "awsuser",
				"DBT_PASSWORD": "my_password",
				"DBT_DB":       "dev",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host:     "examplecluster.abc123xyz789.us-west-1.redshift.amazonaws.com",
						fieldname.Port:     "5439",
						fieldname.User:     "awsuser",
						fieldname.Password: "my_password",
						fieldname.Database: "dev",
					},
				},
			},
		},
	})
}
