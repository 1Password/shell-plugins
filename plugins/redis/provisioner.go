package redis

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type redisArgsProvisioner struct {
}

func redisProvisioner() sdk.Provisioner {
	return redisArgsProvisioner{}
}

func (p redisArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if value, ok := in.ItemFields[fieldname.Password]; ok {
		out.AddEnvVar("REDISCLI_AUTH", value)
	}
	for _, arg := range argsToProvision {
		argName := strings.Split(arg, " ")[0]
		fieldName := sdk.FieldName(strings.Split(arg, " ")[1])
		if fieldValue, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsAtIndex(1, argName, fieldValue)
		}
	}
}

func (p redisArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p redisArgsProvisioner) Description() string {
	return "Provision redis secrets as command-line arguments."
}
