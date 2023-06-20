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
		DocsURL:       sdk.URL("https://cratedb.com/docs/database_credentials"), // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.cratedb.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "CrateDB host to connect to.",
				Optional:            true,
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to CrateDB.",
				Optional:            true,
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "CrateDB user to authenticate as.",
				Optional:            true,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to CrateDB.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryCrateDBConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CRATEPW": fieldname.Password, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Database Credentials in a local config file, and if so,
// implement the function below to add support for importing it.
func TryCrateDBConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config. == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.: config.,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	 string
// }
