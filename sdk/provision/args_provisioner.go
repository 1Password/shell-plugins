package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// ArgsProvisioner provisions secrets as command-line arguments.
// ArgsPosition field is to control where the arguments are placed in the command line. Setting it to 1 provisions the arguments immediately after the executable name, and setting it to the length of the command line provisions the arguments at the end.
type ArgsProvisioner struct {
	sdk.Provisioner

	Schema       map[string]sdk.FieldName
	ArgsPosition map[string]uint
}

// Args creates an ArgsProvisioner that provisions secrets as command line arguments, based
// on the specified schema of field name and argument name, and the position of the argument.
func Args(schema map[string]sdk.FieldName, argsPosition map[string]uint) sdk.Provisioner {
	return ArgsProvisioner{
		Schema:       schema,
		ArgsPosition: argsPosition,
	}
}

func (p ArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			// This safeguard is to ensure that the argsPosition is not out of bounds.
			//
			// Example: For a "redis-cli incr key" command, arguments cannot be provisioned at index 4 or higher.
			if uint(p.ArgsPosition[argName]) > uint(len(out.CommandLine)) {
				out.AddArgsAtIndex(uint(len(out.CommandLine)), argName, value)
				return
			}
			out.AddArgsAtIndex(uint(p.ArgsPosition[argName]), argName, value)
		}
	}
}

func (p ArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p ArgsProvisioner) Description() string {
	return "Provision credentials using command-line args."
}
