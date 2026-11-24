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
	keepProfileWhenYielding(profile, out)
}

// keepProfileWhenYielding passes the profile selected with `--profile` on to the executable as
// AWS_PROFILE when the provisioner supplied no credentials for it.
//
// The `--profile` flag is always stripped from the command line, because an explicit flag makes
// the AWS CLI ignore the provisioned environment credentials. When the provisioner yields (the
// profile belongs to the other credential type), dropping the flag would make the executable run
// against the default profile unless another linked credential provisions for it. Setting
// AWS_PROFILE keeps the user's selection instead: if another linked credential does provision,
// its AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY take precedence over AWS_PROFILE; if none does,
// the executable resolves the requested profile on its own, as it would without the plugin.
func keepProfileWhenYielding(profile string, out *sdk.ProvisionOutput) {
	if profile == "" || len(out.Diagnostics.Errors) > 0 {
		return
	}
	if _, provisioned := out.Environment["AWS_ACCESS_KEY_ID"]; provisioned {
		return
	}
	out.AddEnvVar("AWS_PROFILE", profile)
}

// stripAndReturnProfileFlag strips all occurrences of the `--profile` flag and returns the last occurrence's value.
func stripAndReturnProfileFlag(args []string) (string, []string, error) {
	var profileValue string
	var newArgs []string
	commandOptionsEnded := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		// if `--profile` is used after `--`, it should not be interpreted as a flag
		if arg == "--" {
			commandOptionsEnded = true
		}
		if !commandOptionsEnded && arg == "--profile" {
			if i+1 == len(args) {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}
			profileValue = args[i+1]
			// skip the value too, so later arguments equal to the profile name are kept
			i++
		} else if !commandOptionsEnded && strings.HasPrefix(arg, "--profile=") {
			profileValue = strings.TrimPrefix(arg, "--profile=")
			if profileValue == "" {
				return "", nil, fmt.Errorf("--profile flag was specified without a value")
			}
		} else {
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

// SSOCLIProvisioner provisions SSO role credentials when the AWS CLI is invoked,
// after stripping any --profile flag from the command line so the AWS CLI does
// not try to assume a role on its own.
type SSOCLIProvisioner struct {
}

func (p SSOCLIProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	profile, editedCommandLine, err := stripAndReturnProfileFlag(out.CommandLine)
	if err != nil {
		out.AddError(err)
		return
	}
	if editedCommandLine != nil {
		out.CommandLine = editedCommandLine
	}
	NewSSOProvisioner(profile).Provision(ctx, in, out)
	keepProfileWhenYielding(profile, out)
}

func (p SSOCLIProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p SSOCLIProvisioner) Description() string {
	return "Provision environment variables with temporary AWS IAM Identity Center role credentials AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN"
}
