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
	if editedCommandLine != nil {
		out.CommandLine = editedCommandLine
	}
	stsProvisioner := NewSTSProvisioner(profile)
	stsProvisioner.Provision(ctx, in, out)
}

// stripAndReturnProfileFlag strips all occurrences of the `--profile` flag and returns the last occurrence's value.
func stripAndReturnProfileFlag(args []string) (string, []string, error) {
	var profileValue string
	var newArgs []string
	commandOptionsEnded := false
	for i, arg := range args {
		// if `--profile` is used after `--`, it should not be interpreted as a flag
		if arg == "--" {
			commandOptionsEnded = true
		}
		if !commandOptionsEnded && arg == "--profile" {
			if i+1 == len(args) {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}
			profileValue = args[i+1]
		} else if !commandOptionsEnded && strings.HasPrefix(arg, "--profile=") {
			profileValue = strings.TrimPrefix(arg, "--profile=")
			if profileValue == "" {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}
		} else if commandOptionsEnded || arg != profileValue {
			// make sure only profile specification args are removed from the command line so that aws cli does not attempt to assume a role by itself.
			newArgs = append(newArgs, arg)
		}
	}
	return profileValue, newArgs, nil
}

func (p CLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p CLIProvisioner) Description() string {
	return "Provision environment variables with master credentials or temporary STS credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
