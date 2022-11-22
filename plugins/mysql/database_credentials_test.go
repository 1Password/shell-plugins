package mysql

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMysqlConfigStartsWithClientGroup(t *testing.T) {
	group := "[client]"
	p := sdk.ProvisionInput{
		ItemFields: map[string]string{
			fieldname.User:     "root",
			fieldname.Password: "root",
		},
	}
	config, _ := mysqlConfig(p)

	assert.True(t, strings.HasPrefix(string(config), group), fmt.Sprintf("should start with \"%s\" group", group))
}

func TestMysqlConfigHasPopulatedValues(t *testing.T) {
	p := sdk.ProvisionInput{
		ItemFields: map[string]string{
			"host":     "localhost",
			"port":     "3306",
			"user":     "root",
			"password": "root",
			"database": "db",
		},
	}

	config, _ := mysqlConfig(p)
	entries := strings.Split(string(config), "\n")

	cases := map[string]struct {
		entryKey string
	}{
		"has host value": {
			entryKey: "host",
		},
		"has port value": {
			entryKey: "port",
		},
		"has user value": {
			entryKey: "user",
		},
		"has password value": {
			entryKey: "password",
		},
		"has database value": {
			entryKey: "database",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			hasPopulatedValue := false
			for _, entry := range entries {
				if strings.Contains(entry, tc.entryKey) {
					hasPopulatedValue = true
				}
			}
			assert.True(t, hasPopulatedValue)
		})
	}
}

func TestMysqlConfigHandleEmptyItemFields(t *testing.T) {
	p := sdk.ProvisionInput{
		ItemFields: map[string]string{},
	}
	_, err := mysqlConfig(p)
	if err != nil {
		assert.Fail(t, "should not throw error if no ItemFields")
	}
}

func TestDatabaseCredentialsImporter(t *testing.T) {
	expectedFields := map[string]string{
		"user":     "root",
		"password": "123456",
		"database": "test",
		"port":     "3306",
		"host":     "localhost",
	}

	plugintest.TestImporter(t, DatabaseCredentials().Importer, map[string]plugintest.ImportCase{
		"config file ~/.mysql.cnf": {
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
	plugintest.TestProvisioner(t, DatabaseCredentials().Provisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[string]string{
				"user":     "user",
				"password": "123456",
				"database": "test",
				"host":     "localhost",
				"port":     "3306",
			},
			CommandLine: []string{"mysql"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"mysql", "--defaults-file", "tmp/my.cnf"},
				Files: map[string]sdk.OutputFile{
					"tmp/my.cnf": {
						Contents: []byte("[client]\nuser=user\npassword=123456\nhost=localhost\nport=3306\ndatabase=test\n"),
					},
				},
			},
		},
	})
}
