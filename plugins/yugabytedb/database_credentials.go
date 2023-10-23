package yugabytedb

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
		Name:          credname.DatabaseCredentials,
		DocsURL:       sdk.URL("https://docs.yugabyte.com/preview/admin/ysqlsh/#connect-to-a-database"),
		ManagementURL: sdk.URL("https://cloud.yugabyte.com/clusters"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Yugabyte host to connect to.",
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to Yugabyte.",
				Optional:            true,
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Yugabyte user to get authenticate.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Yugabyte.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to Yugabyte.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"PGHOST":     fieldname.Host,
	"PGPORT":     fieldname.Port,
	"PGUSER":     fieldname.Username,
	"PGPASSWORD": fieldname.Password,
	"PGDATABASE": fieldname.Database,
}
