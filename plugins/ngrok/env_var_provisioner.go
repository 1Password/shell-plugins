package ngrok

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"golang.org/x/mod/semver"
)

type ngrokEnvVarProvisioner struct {
}

func (p ngrokEnvVarProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	currentVersion, err := getNgrokVersion()
	if err != nil {
		out.AddError(err)
		return
	}

	// If the current ngrok CLI version is 3.2.1 or higher,
	// use environment variables to provision the Shell Plugin credentials
	//
	// semver.Compare resulting in 0 means 3.2.1 is in use
	// semver.Compare resulting in +1 means >3.2.1 is in use
	if semver.Compare(currentVersion, envVarAuthVersion) == -1 {
		out.AddError(fmt.Errorf("ngrok version %s is not supported. Please upgrade to version %s or higher", currentVersion, envVarAuthVersion))
		return
	}

	out.AddEnvVar("NGROK_AUTHTOKEN", in.ItemFields[fieldname.Authtoken])
	out.AddEnvVar("NGROK_API_KEY", in.ItemFields[fieldname.APIKey])
}

func (p ngrokEnvVarProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p ngrokEnvVarProvisioner) Description() string {
	return "Provision ngrok credentials as environment variables NGROK_AUTH_TOKEN and NGROK_API_KEY"
}
