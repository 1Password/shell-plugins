package mongodb

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, mongodbShellProvisioner(), map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Host:     "localhost",
				fieldname.Port:     "27017",
				fieldname.Username: "default",
				fieldname.Password: "password",
			},
			CommandLine: []string{"mongosh"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"mongosh", "--password", "password", "--username", "default", "--port", "27017", "--host", "localhost"}, // Each argument is provisioned at index 1, pushing the existing arguments forward to the next index
			},
		},
	})
}
