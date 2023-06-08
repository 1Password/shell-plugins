package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// ArgsProvisioner provisions secrets as command-line arguments.
// ArgsPosition field is to control where the arguments are placed in the command line. A value of "true" provisions the arguments immediately after the executable name, or at the end of the command.
type ArgsProvisioner struct {
	sdk.Provisioner

	Schema       map[string]sdk.FieldName
	ArgsPosition map[string]bool
}

// Args creates an ArgsProvisioner that provisions secrets as command line arguments, based
// on the specified schema of field name and argument name, and the value of ArgsPosition.
func Args(schema map[string]sdk.FieldName, argsPosition map[string]bool) sdk.Provisioner {
	return ArgsProvisioner{
		Schema:       schema,
		ArgsPosition: argsPosition,
	}
}

func (p ArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgs(p.ArgsPosition[argName], argName, value)
		}
	}
}

func (p ArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p ArgsProvisioner) Description() string {
	return "Provision credentials using command-line args."
}
