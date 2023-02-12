package sapbtp

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Username: "john.doe@email.com",
				fieldname.Password: "CpfSh78WLKpO6teQitmYbH7EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{
					"--user",
					"john.doe@email.com",
					"--password",
					"CpfSh78WLKpO6teQitmYbH7EXAMPLE",
				},
			},
		},
	})
}
