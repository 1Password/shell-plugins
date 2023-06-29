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
				fieldname.User: "vertica",
				fieldname.Password: "",
				fieldname.Host: "localhost",
				fieldname.Port: "5433",
				fieldname.Database: "VMart",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"VSQL_USER": "vertica", 
					"VSQL_PASSWORD": "", 
					"VSQL_HOST": "localhost", 
					"VSQL_PORT": "5433", 
					"VSQL_DATABASE": "VMart",
				},
			},
		},
	})
}
