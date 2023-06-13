package mongodbshell

import (
	"context"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

type mongodbShellArgsProvisioner struct {
}

func mongodbShellProvisioner() sdk.Provisioner {
	return mongodbShellArgsProvisioner{}
}

func (p mongodbShellArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for _, arg := range argsToProvision {
		argName := strings.Split(arg, " ")[0]
		fieldName := sdk.FieldName(strings.Split(arg, " ")[1])
		if fieldValue, ok := in.ItemFields[fieldName]; ok {
			out.AddArgsAtIndex(1, argName, fieldValue)
		}
	}
}

func (p mongodbShellArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p mongodbShellArgsProvisioner) Description() string {
	return "Provision MongoDB Shell secrets as command-line arguments."
}
