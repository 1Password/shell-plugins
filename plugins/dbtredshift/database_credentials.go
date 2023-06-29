package dbtredshift

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
		DocsURL: sdk.URL("https://docs.getdbt.com/docs/core/connect-data-platform/redshift-setup#password-based-authentication"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Redshift host to connect to.",
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to Redshift.",
				Optional:            true,
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "Redshift user to authenticate as.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Redshift.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to.",
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DBT_HOST":     fieldname.Host,
	"DBT_PORT":     fieldname.Port,
	"DBT_USER":     fieldname.User,
	"DBT_PASSWORD": fieldname.Password,
	"DBT_DB":       fieldname.Database,
}
