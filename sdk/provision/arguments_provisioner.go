package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// ArgumentsProvisioner provisions secrets as command line arguments
type ArgumentsProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

// Arguments creates an ArgumentsProvisioner that provisions secrets as command line arguments, based
// on the specified schema of argument flag/name and its value, which is a FieldName.
func Arguments(schema map[string]sdk.FieldName) sdk.Provisioner {
	return ArgumentsProvisioner{
		Schema: schema,
	}
}

func (p ArgumentsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgs(argName, value)
		}
	}
}

func (p ArgumentsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: passed argument values get wiped automatically when the process exits.
}

func (p ArgumentsProvisioner) Description() string {
	return "Provision secrets as command line arguments"
}
