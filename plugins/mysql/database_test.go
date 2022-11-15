package mysql

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
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
		"has Host value": {
			entryKey: strings.ToLower(fieldname.Host),
		},
		"has Port value": {
			entryKey: strings.ToLower(fieldname.Port),
		},
		"has User value": {
			entryKey: strings.ToLower(fieldname.User),
		},
		"has Password value": {
			entryKey: strings.ToLower(fieldname.Password),
		},
		"has Database value": {
			entryKey: strings.ToLower(fieldname.Database),
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
