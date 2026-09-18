package ssh

import (
	"context"
	"github.com/1Password/shell-plugins/sdk/importer"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func UserLogin() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.UserLogin,
		DocsURL: sdk.URL("https://linux.die.net/man/1/sshpass"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				AlternativeNames:    []string{fieldname.HostAddress.String(), fieldname.URL.String()},
				MarkdownDescription: "SSH Host to connect to.",
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "User to authenticate as.",
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate.",
				Secret:              true,
			},
		},
		DefaultProvisioner: sshProvisioner{},
		Importer:           importer.NoOp(),
	}
}

type sshProvisioner struct{}

func (s sshProvisioner) Description() string {
	return "SSH password provisioner"
}

func (s sshProvisioner) Provision(_ context.Context, input sdk.ProvisionInput, output *sdk.ProvisionOutput) {
	output.AddArgs("-p", input.ItemFields[fieldname.Password],
		"ssh", input.ItemFields[fieldname.Username]+"@"+input.ItemFields[fieldname.Host])
}

func (s sshProvisioner) Deprovision(_ context.Context, _ sdk.DeprovisionInput, _ *sdk.DeprovisionOutput) {
	// no-op
}
