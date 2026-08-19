package ssh

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestUserLoginProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, UserLogin().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Password: "zhexample",
				fieldname.Username: "user",
				fieldname.Host:     "192.168.1.20",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"-p", "zhexample", "ssh", "user@192.168.1.20"},
			},
		},
	})
}
