package cratedb

import (
	"context"

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
		DocsURL:       sdk.URL("https://crate.io/docs/crate/crash/en/latest/run.html#environment-variables"),
		ManagementURL: sdk.URL("https://console.cratedb.cloud/account/settings"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "CrateDB host to connect to.",
				Optional:            true,
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "CrateDB user to authenticate as.",
				Optional:            true,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to CrateDB.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CRATEPW": fieldname.Password, 
}

