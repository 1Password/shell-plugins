package provision

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

// EnvVarProvisioner provisions secrets as environment variables.
type EnvVarProvisioner struct {
	sdk.Provisioner

	Schema map[string]sdk.FieldName
}

// EnvVars creates an EnvVarProvisioner that provisions secrets as environment variables, based
// on the specified schema of field name and environment variable name.
func EnvVars(schema map[string]sdk.FieldName) sdk.Provisioner {
	return EnvVarProvisioner{
		Schema: schema,
	}
}

func (p EnvVarProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for envVarName, fieldName := range p.Schema {
		if value, ok := in.ItemFields[fieldName]; ok {
			out.AddEnvVar(envVarName, value)
		}
	}
}

func (p EnvVarProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p EnvVarProvisioner) Description() string {
	var envVarNames []string
	for envVarName := range p.Schema {
		envVarNames = append(envVarNames, envVarName)
	}

	return fmt.Sprintf("Provision environment variables: %s", strings.Join(envVarNames, ", "))
}
