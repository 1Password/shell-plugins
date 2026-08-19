package cockroachdb

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
		DocsURL:       sdk.URL("https://www.cockroachlabs.com/docs/stable/connection-parameters.html"),
		ManagementURL: sdk.URL("https://cockroachlabs.cloud/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "CockroachDB host to connect to.",
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to CockroachDB.",
				Optional:            true,
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "CockroachDB user to authenticate as.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to CockroachDB.",
				Secret:              true,
				Optional:            true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to. Defaults to 'defaultdb'.",
				Optional:            true,
			},
			{
				Name:                "insecure",
				MarkdownDescription: "Connect in insecure mode (skip TLS verification). Set to '1' to skip TLS verification.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"COCKROACH_HOST":     fieldname.Host,
	"COCKROACH_PORT":     fieldname.Port,
	"COCKROACH_USER":     fieldname.User,
	"COCKROACH_PASSWORD": fieldname.Password,
	"COCKROACH_DATABASE": fieldname.Database,
	"COCKROACH_INSECURE": "insecure",
}
