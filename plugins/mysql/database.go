package mysql

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
		DocsURL:       sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/connecting.html"),
		ManagementURL: sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/mysql-config-editor.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Credentials,
				MarkdownDescription: "Credentials used to authenticate to MySQL.",
				Secret:              true,
			},
		},
		Provisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryMySQLConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.Credentials: "MYSQL_DATABASE_CREDENTIALS", // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Database Credentials in a local config file, and if so,
// implement the function below to add support for importing it.
func TryMySQLConfigFile() sdk.Importer {
	return importer.TryFile("~/path/to/config/file.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportOutput) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Credentials == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: []sdk.ImportCandidateField{
		// 		{
		// 			Field: fieldname.Credentials,
		// 			Value: config.Credentials,
		// 		},
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Credentials string
// }
