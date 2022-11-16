package mysql

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"os"
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
			fieldname.Host:     "localhost",
			fieldname.Port:     "3306",
			fieldname.User:     "root",
			fieldname.Password: "root",
			fieldname.Database: "db",
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

func TestTryCredentialsFile(t *testing.T) {
	config, _ := mysqlConfig(sdk.ProvisionInput{
		ItemFields: map[string]string{
			fieldname.Host:     "localhost",
			fieldname.Port:     "3306",
			fieldname.User:     "root",
			fieldname.Password: "root",
			fieldname.Database: "db",
		},
	})
	path, _ := os.Getwd()
	configFilePath := fmt.Sprintf("%s/mysql-cred.cnf", path)
	os.WriteFile(configFilePath, config, 0644)

	res := TryMySQLConfigFile(configFilePath)
	out := &sdk.ImportOutput{}
	res(context.TODO(), sdk.ImportInput{}, out)

	candidates := out.AllCandidates()
	assert.True(t, len(candidates) == 1)

	cases := map[string]struct {
		entryKey string
	}{
		"ImportCandidate has host value": {
			entryKey: "host",
		},
		"ImportCandidate has port value": {
			entryKey: "port",
		},
		"ImportCandidate has user value": {
			entryKey: "user",
		},
		"ImportCandidate has password value": {
			entryKey: "password",
		},
		"ImportCandidate has database value": {
			entryKey: "database",
		},
	}

	candidate := candidates[0]
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			hasPopulatedValue := false
			for _, f := range candidate.Fields {
				if strings.Contains(f.Field, tc.entryKey) {
					hasPopulatedValue = true
				}
			}
			assert.True(t, hasPopulatedValue)
		})
	}

	os.Remove(configFilePath)
}
