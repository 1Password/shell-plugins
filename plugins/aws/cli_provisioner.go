package aws

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
)

type awsCLIProvisioner struct {
}

func NewAwsCLIProvisioner() sdk.Provisioner {
	return awsCLIProvisioner{}
}

func (p awsCLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile, err := stripAndReturnProfileFlag(&out.CommandLine)
	if err != nil {
		out.AddError(err)
		return
	}

	NewSTSProvisioner(profile).Provision(ctx, in, out)
}

func stripAndReturnProfileFlag(args *[]string) (string, error) {
	for i, arg := range *args {
		if arg == "--profile" {
			if i+1 == len(*args) {
				return "", fmt.Errorf("--profile flag was specified without a value")

			}
			profile := (*args)[i+1]
			// Remove --profile flag so aws cli doesn't attempt to assume role by itself
			*args = removeProfileFlagFromArgs(i, *args)
			return profile, nil
		}
	}
	return "", nil
}

func removeProfileFlagFromArgs(argIndex int, args []string) []string {
	result := append(args[0:argIndex], args[argIndex+2:]...)
	return result
}

func (p awsCLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p awsCLIProvisioner) Description() string {
	return "Provision environment variables with master credentials or temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
