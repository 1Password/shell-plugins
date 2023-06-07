package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

// ChainedProvisioner chains multiple provisioners together for use at the same time.
type ChainedProvisioner struct {
	sdk.Provisioner

	Provisioners []sdk.Provisioner
}

// ChainProvisioners creates a ChainedProvisioner that chains multiple provisioners together for use at the same time.
func ChainProvisioners(provisioners ...sdk.Provisioner) sdk.Provisioner {
	return ChainedProvisioner{
		Provisioners: provisioners,
	}
}

func (p ChainedProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for _, provisioner := range p.Provisioners {
		provisioner.Provision(ctx, in, out)
	}
}

func (p ChainedProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables and args get wiped automatically when the process exits.
}

func (p ChainedProvisioner) Description() string {
	return "Handle multiple provisioners at once"
}
