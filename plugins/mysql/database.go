package mysql

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"strings"
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.DatabaseCredentials,
		DocsURL:       sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/connecting.html"),
		ManagementURL: sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/mysql-config-editor.html"),
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
		Provisioner: provision.TempFile(mysqlConfig, provision.SetPathAsArg("--defaults-file")),
	}
}

func mysqlConfig(in sdk.ProvisionInput) ([]byte, error) {
	content := "[client]" + "\n"
	configSchema := map[string]string{
		fieldname.User:     "",
		fieldname.Password: "",
		fieldname.Host:     "127.0.0.1", // Default host
		fieldname.Port:     "3306",      // Default port
		fieldname.Database: "",
	}

	for key, defaultValue := range configSchema {
		configEntryVal := in.ItemFields[key]
		if configEntryVal == "" {
			configEntryVal = defaultValue
		}

		if configEntryVal != "" {
			configEntry := []string{strings.ToLower(key), "=", configEntryVal, "\n"}
			content += strings.Join(configEntry, "")
		}
	}

	return []byte(content), nil
}
