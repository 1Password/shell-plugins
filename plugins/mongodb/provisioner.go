package mongodb

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

var argsToProvision = []string{
	"--host", fieldname.Host.String(),
	"--port", fieldname.Port.String(),
	"--username", fieldname.Username.String(),
	"--password", fieldname.Password.String(),
}

type mongodbShellArgsProvisioner struct {
}

func mongodbShellProvisioner() sdk.Provisioner {
	return mongodbShellArgsProvisioner{}
}

func (p mongodbShellArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for i, arg := range argsToProvision {
		if i%2 == 0 {
			argName := arg
			fieldName := sdk.FieldName(argsToProvision[i+1])
			if fieldValue, ok := in.ItemFields[fieldName]; ok {
				out.AddArgsAtIndex(1, argName, fieldValue)
			}
		}
	}
}

func (p mongodbShellArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p mongodbShellArgsProvisioner) Description() string {
	return "Provision MongoDB Shell secrets as command-line arguments."
}
