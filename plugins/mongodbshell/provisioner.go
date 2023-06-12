package mongodbshell

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type mongodbShellArgsProvisioner struct {
	Schema map[string]sdk.FieldName
}

func mongodbShellFlags(schema map[string]sdk.FieldName) sdk.Provisioner {
	return mongodbShellArgsProvisioner{
		Schema: schema,
	}
}

func (p mongodbShellArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsAtIndex(1, argName, value)
		}
	}
}

func (p mongodbShellArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p mongodbShellArgsProvisioner) Description() string {
	return "Provision MongoDB Shell secrets as command-line arguments."
}
