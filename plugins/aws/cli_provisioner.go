package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
)

type CLIProvisioner struct {
}

func (p CLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile, editedCommandLine, err := stripAndReturnProfileFlag(out.CommandLine)
	if err != nil {
		out.AddError(err)
		return
	}
	out.CommandLine = editedCommandLine
	stsProvisioner := STSProvisioner{profileName: profile}
	stsProvisioner.Provision(ctx, in, out)
}

func stripAndReturnProfileFlag(args []string) (string, []string, error) {
	for i, arg := range args {
		// if `--profile` is used after `--`, it should not be interpreted as a flag
		if arg == "--" {
			break
		}

		if arg == "--profile" {
			if i+1 == len(args) {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}
			profile := args[i+1]

			// Remove --profile flag so aws cli doesn't attempt to assume role by itself
			args = append(args[0:i], args[i+2:]...)

			return profile, args, nil
		} else if strings.HasPrefix(arg, "--profile=") {
			profile := strings.TrimPrefix(arg, "--profile=")
			if profile == "" {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}

			// Remove --profile flag so aws cli doesn't attempt to assume role by itself
			args = append(args[0:i], args[i+1:]...)

			return profile, args, nil
		}
	}
	return "", nil, nil
}

func (p CLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p CLIProvisioner) Description() string {
	return "Provision environment variables with master credentials or temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
