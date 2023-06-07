package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type ArgsProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

func Args(schema map[string]sdk.FieldName) sdk.Provisioner {
	return ArgsProvisioner{
		Schema: schema,
	}
}

func (p ArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsFirst(argName, value)
		}
	}
}

func (p ArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p ArgsProvisioner) Description() string {
	return "Provision credentials using command-line args."
}
