package influxdb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func InfluxDBCLI() schema.Executable {
	return schema.Executable{
		Name:    "InfluxDB CLI",
		Runs:    []string{"influx"},
		DocsURL: sdk.URL("https://docs.influxdata.com/influxdb/cloud/tools/influx-cli/"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.NotWhenContainsArgs("config"),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.DatabaseCredentials,
			},
		},
	}
}
