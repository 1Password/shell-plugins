package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// EnvVarProvisioner provisions secrets as command-line arguments.
// ArgsPosition field is to control where the arguments are placed in the command line. We check for only "first" and provision those arguments first, otherwise we provision them last.
type ArgsProvisioner struct {
	sdk.Provisioner

	Schema       map[string]sdk.FieldName
	ArgsPosition string
}

// Args creates an ArgsProvisioner that provisions secrets as command line arguments, based
// on the specified schema of field name and argument name.
func Args(schema map[string]sdk.FieldName, argsPosition string) sdk.Provisioner {
	return ArgsProvisioner{
		Schema:       schema,
		ArgsPosition: argsPosition,
	}
}

func (p ArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if p.ArgsPosition == "first" {
		for argName, fieldName := range p.Schema {
			if value, ok := in.ItemFields[fieldName]; ok {
				out.AddArgsFirst(argName, value)
			}
		}
		return
	}

	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsLast(argName, value)
		}
	}
}

func (p ArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p ArgsProvisioner) Description() string {
	return "Provision credentials using command-line args."
}
