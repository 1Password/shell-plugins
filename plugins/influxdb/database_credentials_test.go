package influxdb

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
				fieldname.Host:         "https://us-west-2-1.aws.cloud2.influxdata.com",
				fieldname.Organization: "1Password.com",
				fieldname.AccessToken:  "BHsmEerxKV2yDaNNv31lPHMEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"INFLUX_HOST":  "https://us-west-2-1.aws.cloud2.influxdata.com",
					"INFLUX_ORG":   "1Password.com",
					"INFLUX_TOKEN": "BHsmEerxKV2yDaNNv31lPHMEXAMPLE",
				},
			},
		},
	})
}
