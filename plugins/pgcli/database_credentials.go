package pgcli

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
		DocsURL: sdk.URL("https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Postgres host to connect to.",
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to Postgres.",
				Optional:            true,
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "Postgres user to authenticate as.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Postgres.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to. Defaults to the name of the authenticated user.",
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
	"PGUSER":     fieldname.User,
	"PGPASSWORD": fieldname.Password,
	"PGDATABASE": fieldname.Database,
}
