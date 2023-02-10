package snowflake

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestLoginDetailsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, LoginDetails().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Account:  "accountname",
				fieldname.Username: "username",
				fieldname.Password: "password1234",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNOWSQL_ACCOUNT": "accountname",
					"SNOWSQL_USER":    "username",
					"SNOWSQL_PWD":     "password1234",
				},
			},
		},
	})
}

func TestLoginDetailsImporter(t *testing.T) {
	expectedFields := map[sdk.FieldName]string{
		fieldname.Account:  "accountname",
		fieldname.Username: "username",
		fieldname.Password: "password1234",
	}

	plugintest.TestImporter(t, LoginDetails().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SNOWSQL_ACCOUNT": "accountname",
				"SNOWSQL_USER":    "username",
				"SNOWSQL_PWD":     "password1234",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
			},
		},
		"config file ([connections] section only)": {
			Files: map[string]string{
				"~/.snowsql/config": plugintest.LoadFixture(t, "config1"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
			},
		},
		"config file ([connections] and [connections.example] sections)": {
			Files: map[string]string{
				"~/.snowsql/config": plugintest.LoadFixture(t, "config2"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
			},
		},
		"config file ([connections.example] section first)": {
			Files: map[string]string{
				"~/.snowsql/config": plugintest.LoadFixture(t, "config3"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{Fields: expectedFields},
			},
		},
	})
}
