package pypi

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// pypiToolProvisioner sets username and password environment variables
// for PyPI publishing tools. The username is always "__token__" when
// using API token authentication.
type pypiToolProvisioner struct {
	usernameEnvVar string
	passwordEnvVar string
}

func PyPIToolProvisioner(usernameEnvVar, passwordEnvVar string) sdk.Provisioner {
	return pypiToolProvisioner{
		usernameEnvVar: usernameEnvVar,
		passwordEnvVar: passwordEnvVar,
	}
}

func (p pypiToolProvisioner) Description() string {
	return fmt.Sprintf("Provision PyPI API token via %s and %s environment variables", p.usernameEnvVar, p.passwordEnvVar)
}

func (p pypiToolProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	out.AddEnvVar(p.usernameEnvVar, "__token__")
	out.AddEnvVar(p.passwordEnvVar, in.ItemFields[fieldname.Token])
}

func (p pypiToolProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Environment variables are automatically cleaned up when the process ends.
}