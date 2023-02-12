package sapbtp

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type btpProvisioner struct {
}

func BTPProvisioner() sdk.Provisioner {
	return btpProvisioner{}
}

func (p btpProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {

	if username, ok := in.ItemFields[fieldname.Username]; ok {
		out.AddArgs("--user", username)
	}

	if password, ok := in.ItemFields[fieldname.Password]; ok {
		out.AddArgs("--password", password)
	}

	return
}

func (p btpProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p btpProvisioner) Description() string {
	return "Provision command line variables --user and --password with BTP cli"
}
