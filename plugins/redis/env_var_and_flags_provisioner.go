package redis

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type EnvVarFlagsProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

func EnvVarFlags(schema map[string]sdk.FieldName) sdk.Provisioner {
	return EnvVarFlagsProvisioner{
		Schema: schema,
	}
}

func (p EnvVarFlagsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if value, ok := in.ItemFields[fieldname.Password]; ok {
		out.AddEnvVar("REDISCLI_AUTH", value)
	}
	for flagName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsImmediatelyAfterExecutableName(flagName, value)
		}
	}
}

func (p EnvVarFlagsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p EnvVarFlagsProvisioner) Description() string {
	return "Provision environment variables with master credentials or temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
