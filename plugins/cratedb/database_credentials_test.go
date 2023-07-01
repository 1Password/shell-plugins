package cratedb

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestDatabaseCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, CrateArgsProvisioner{}, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ 
				fieldname.Password: "1<34&f0rg3t@me",
				fieldname.Host: "https://love.aks1.eastus2.azure.cratedb.net:4200",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CRATEPW": "1<34&f0rg3t@me",
				},
				CommandLine: []string{
					"--username",
					"",
					"--hosts",
					"https://love.aks1.eastus2.azure.cratedb.net:4200",
				},
			},
		},
	})
}
