package cratedb

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type crateArgsProvisioner struct {
}

func crateProvisioner() sdk.Provisioner {
	return crateArgsProvisioner{}
}

func (p crateArgsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if value, ok := in.ItemFields[fieldname.Password]; ok {
		out.AddEnvVar("CRATEPW", value)
	}
	
	var user, host string
	if fieldValue, ok := in.ItemFields[fieldname.Username]; ok {
		user=fieldValue
	}
	if fieldValue, ok := in.ItemFields[fieldname.Host]; ok {
		host=fieldValue
	}
		commandLine := []string{out.CommandLine[0], "--username", user, "--hosts", host, }
		commandLine = append(commandLine, out.CommandLine[1:]...)
		out.CommandLine = commandLine
	
}


func (p crateArgsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: credentials get wiped automatically when the process exits.
}

func (p crateArgsProvisioner) Description() string {
	return "Provision CrateDB username, host as command-line arguments && Password as Env ."
}
