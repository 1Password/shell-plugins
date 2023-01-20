package snowflake

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
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

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SNOWSQL_ACCOUNT": "accountname",
				"SNOWSQL_USER":    "username",
				"SNOWSQL_PWD":     "password1234",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:  "accountname",
						fieldname.Username: "username",
						fieldname.Password: "password1234",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.snowsql/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Account:  "accountname",
						fieldname.Username: "username",
						fieldname.Password: "password1234",
					},
				},
			},
		},
	})
}
