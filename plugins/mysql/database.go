package mysql

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.DatabaseCredentials,
		DocsURL: sdk.URL("https://dev.mysql.com/doc/refman/8.0/en/connecting.html"),
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

type configValue struct {
	value        string
	defaultValue string
}

type configKey string

const (
	host     configKey = "host"
	port     configKey = "port"
	user     configKey = "user"
	password configKey = "password"
	database configKey = "database"
)

func mysqlConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := map[string]string{
		"host": "127.0.0.1", // Default host
		"port": "3306",      // Default port
	}

	if user, ok := in.ItemFields[fieldname.User]; ok {
		config["user"] = user
	}

	if password, ok := in.ItemFields[fieldname.Password]; ok {
		config["password"] = password
	}

	if host, ok := in.ItemFields[fieldname.Host]; ok {
		config["host"] = host
	}

	if port, ok := in.ItemFields[fieldname.Port]; ok {
		config["port"] = port
	}

	if database, ok := in.ItemFields[fieldname.Database]; ok {
		config["database"] = database
	}

	content := "[client]\n"
	for key, value := range config {
		configLine := fmt.Sprintf("%s=%s\n", key, value)
		content += configLine
	}

	return []byte(content), nil
}
