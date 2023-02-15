package mysql

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseCredentialsImporter(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.User:     "root",
		fieldname.Password: "123456",
		fieldname.Database: "test",
		fieldname.Port:     "3306",
		fieldname.Host:     "localhost",
	}

	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"MySQL config file": {
			Files: map[string]string{
				"/etc/my.cnf":       plugintest.LoadFixture(t, "mysql.cnf"),
				"/etc/mysql/my.cnf": plugintest.LoadFixture(t, "mysql.cnf"),
				"~/.my.cnf":         plugintest.LoadFixture(t, "mysql.cnf"),
				"~/.mylogin.cnf":    plugintest.LoadFixture(t, "mysql.cnf"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
				{Fields: expectedFields},
				{Fields: expectedFields},
				{Fields: expectedFields},
			},
		},
	})
}

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.User:     "root",
				fieldname.Password: "123456",
				fieldname.Database: "test",
				fieldname.Host:     "localhost",
				fieldname.Port:     "3306",
			},
			CommandLine: []string{"mysql"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"mysql", "--defaults-file=/tmp/my.cnf"},
				Files: map[string]sdk.OutputFile{
					"/tmp/my.cnf": {
						Contents: []byte(plugintest.LoadFixture(t, "mysql.cnf")),
					},
				},
			},
		},
	})
}

func TestMysqlConfigHandleEmptyItemFields(t *testing.T) {
	p := sdk.ProvisionInput{
		ItemFields: map[sdk.FieldName]string{},
	}
	_, err := mysqlConfig(p, nil)
	if err != nil {
		assert.Fail(t, "should not throw error if no ItemFields")
	}
}
