package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type FlagsProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

func Flags(schema map[string]sdk.FieldName) sdk.Provisioner {
	return FlagsProvisioner{
		Schema: schema,
	}
}

func (p FlagsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for flagName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsFirst(flagName, value)
		}
	}
}

func (p FlagsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p FlagsProvisioner) Description() string {
	return "Provision credentials using command-line flags."
}
