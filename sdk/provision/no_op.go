package provision

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
)

type noOp struct {
	sdk.Provisioner
}

// NoOp can be used as a provisioner stub while developing plugins.
func NoOp() sdk.Provisioner {
	return noOp{}
}

func (p noOp) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	// No op
}

func (p noOp) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// No op
}

func (p noOp) Info() sdk.ProvisionerInfo {
	return sdk.ProvisionerInfo{
		Description: "No op",
	}
}
