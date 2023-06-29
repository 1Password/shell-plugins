package influxdb

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.DatabaseCredentials,
		DocsURL: sdk.URL("https://docs.influxdata.com/influxdb/v2.7/reference/cli/influx/config/create/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "InfluxDB host name",
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "InfluxDB Organization name",
			},
			{
				Name:                fieldname.AccessToken,
				MarkdownDescription: "InfluxDB Token value",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"INFLUX_HOST":  fieldname.Host,
	"INFLUX_ORG":   fieldname.Organization,
	"INFLUX_TOKEN": fieldname.AccessToken,
}
