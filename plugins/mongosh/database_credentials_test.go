package mongosh

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestUserLoginProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DatabaseCredentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			CommandLine: []string{"mongosh", "test.js"},
			ItemFields: map[sdk.FieldName]string{
				fieldname.Username: "aexample",
				fieldname.Password: "apassword",
				fieldname.Host:     "example.org",
				fieldname.Port:     "2121",
				fieldname.Database: "example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"mongosh", "--username", "aexample", "--password", "apassword", "--host", "example.org", "--port", "2121", "example", "test.js"},
			},
		},
	})
}
