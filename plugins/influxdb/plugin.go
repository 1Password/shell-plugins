package influxdb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "influxdb",
		Platform: schema.PlatformInfo{
			Name:     "InfluxDB",
			Homepage: sdk.URL("https://docs.influxdata.com/influxdb/"),
		},
		Credentials: []schema.CredentialType{
			DatabaseCredentials(),
		},
		Executables: []schema.Executable{
			InfluxDBCLI(),
		},
	}
}
