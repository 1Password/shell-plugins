package mongosh

import (
	"context"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func DatabaseCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.DatabaseCredentials,
		DocsURL: sdk.URL("https://www.mongodb.com/docs/mongodb-shell/connect/#connect-with-ldap"),
		//ManagementURL: sdk.URL("https://console.mongosh.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "MongoDB host to connect to.",
				Optional:            true,
			},
			{
				Name:                fieldname.Port,
				MarkdownDescription: "Port used to connect to MongoDB.",
				Optional:            true,
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "MongoDB user to authenticate as.",
				Optional:            true,
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to MongoDB.",
				Secret:              true,
			},
			{
				Name:                fieldname.Database,
				MarkdownDescription: "Database name or URL to connect to.",
				Optional:            true,
			},
		},
		DefaultProvisioner: mongoshProvisioner{},
	}
}

type mongoshProvisioner struct{}

func (m mongoshProvisioner) Description() string {
	return "mongosh cli provisioner"
}

func (m mongoshProvisioner) Provision(ctx context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	if username, ok := input.ItemFields[fieldname.Username]; ok {
		output.AddArgs("--username", username)
	}

	if password, ok := input.ItemFields[fieldname.Password]; ok {
		output.AddArgs("--password", password)
	}

	if host, ok := input.ItemFields[fieldname.Host]; ok {
		output.AddArgs("--host", host)
	}

	if port, ok := input.ItemFields[fieldname.Port]; ok {
		output.AddArgs("--port", port)
	}

	if db, ok := input.ItemFields[fieldname.Database]; ok {
		output.AddArgs(db)
	}
}

func (m mongoshProvisioner) Deprovision(ctx context.Context, input sdk.DeprovisionInput, output *sdk.DeprovisionOutput) {
	// No-op
}
