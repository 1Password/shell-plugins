package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/1Password/shell-plugins/sdk"
)

const defaultProfileName = "default"

type awsCLIProvisioner struct {
}

func newAwsCLIProvisioner() sdk.Provisioner {
	return awsCLIProvisioner{}
}

func (p awsCLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profileName := defaultProfileName

	if profile := stripAndReturnProfileFlag(out); profile != nil {
		profileName = *profile
	} else if profile := os.Getenv("AWS_PROFILE"); profile != "" {
		profileName = profile
	}
	newSTSProvisioner(profileName).Provision(ctx, in, out)
}

func stripAndReturnProfileFlag(out *sdk.ProvisionOutput) *string {
	for i, arg := range out.CommandLine {
		if arg == "--profile" {
			if i+1 == len(out.CommandLine) {
				out.AddError(fmt.Errorf("--profile flag was specified without a value"))
				return nil
			}
			profile := out.CommandLine[i+1]
			// Remove --profile flag so aws cli doesn't attempt to assume role by itself
			out.CommandLine = removeProfileFlagFromArgs(i, out.CommandLine)
			return &profile
		}
	}
	return nil
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
