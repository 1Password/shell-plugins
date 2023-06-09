package redis

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type redisArgsProvisioner struct {
	AddArgsAtIndex map[string]uint
	Schema         map[string]sdk.FieldName
}

func redisFlags(addArgsAtIndex map[string]uint, schema map[string]sdk.FieldName) sdk.Provisioner {
	return redisArgsProvisioner{
		AddArgsAtIndex: addArgsAtIndex,
		Schema:         schema,
	}
}

func (p redisArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for argName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsAtIndex(p.AddArgsAtIndex[argName], argName, value)
		}
	}
}

func (p redisArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p redisArgsProvisioner) Description() string {
	return "Provision redis secrets as command-line arguments."
}
