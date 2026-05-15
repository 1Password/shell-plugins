package vertica

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
		DocsURL: sdk.URL("https://www.vertica.com/docs/9.2.x/HTML/Content/Authoring/AdministratorsGuide/DBUsersAndPrivileges/Users/CreatingADatabaseUser.htm"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "Vertica host to connect to.",
				Optional:            true,
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to Vertica.",
				Optional:            true,
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Vertica user to authenticate as.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Vertica.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to.",
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"VSQL_USER":     fieldname.Username,
	"VSQL_PASSWORD": fieldname.Password,
	"VSQL_HOST":     fieldname.Host,
	"VSQL_PORT":     fieldname.Port,
	"VSQL_DATABASE": fieldname.Database,
}
