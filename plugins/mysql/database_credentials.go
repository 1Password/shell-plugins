package mysql

import (
	"context"
	"fmt"

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
		DocsURL: sdk.URL("https://dev.mysql.com/doc/refman/en/connecting.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "MySQL host to connect to.",
				Optional:            true,
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to MySQL.",
				Optional:            true,
			},
			{
				Name:                fieldname.User,
				MarkdownDescription: "MySQL user to authenticate as.",
				Optional:            true,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to MySQL.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name to connect to.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.TempFile(mysqlConfig, provision.Filename("my.cnf"), provision.AddArgs("--defaults-file={{ .Path }}")),
		Importer: importer.TryAll(
			TryMySQLConfigFile("/etc/my.cnf"),
			TryMySQLConfigFile("/etc/mysql/my.cnf"),
			TryMySQLConfigFile("~/.my.cnf"),
			TryMySQLConfigFile("~/.mylogin.cnf"),
		),
	}
}

func mysqlConfig(in sdk.ProvisionInput) ([]byte, error) {
	content := "[client]\n"

	if user, ok := in.ItemFields[fieldname.User]; ok {
		content += configFileEntry("user", user)
	}

	if password, ok := in.ItemFields[fieldname.Password]; ok {
		content += configFileEntry("password", password)
	}

	if host, ok := in.ItemFields[fieldname.Host]; ok {
		content += configFileEntry("host", host)
	}

	if port, ok := in.ItemFields[fieldname.Port]; ok {
		content += configFileEntry("port", port)
	}

	if database, ok := in.ItemFields[fieldname.Database]; ok {
		content += configFileEntry("database", database)
	}

	return []byte(content), nil
}

func TryMySQLConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("user") && section.Key("user").Value() != "" {
				fields[fieldname.User] = section.Key("user").Value()
			}

			if section.HasKey("password") && section.Key("password").Value() != "" {
				fields[fieldname.Password] = section.Key("password").Value()
			}

			if section.HasKey("database") && section.Key("database").Value() != "" {
				fields[fieldname.Database] = section.Key("database").Value()
			}

			if section.HasKey("host") && section.Key("host").Value() != "" {
				fields[fieldname.Host] = section.Key("host").Value()
			}

			if section.HasKey("port") && section.Key("port").Value() != "" {
				fields[fieldname.Port] = section.Key("port").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}

func configFileEntry(key string, value string) string {
	return fmt.Sprintf("%s=%s\n", key, value)
}
