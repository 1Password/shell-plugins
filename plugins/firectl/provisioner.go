package firectl

import (
	"context"
	"errors"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type firectlCLIArgProvisioner struct {
	sdk.Provisioner
}

func FirectlCLIArgProvisioner() sdk.Provisioner {
	return firectlCLIArgProvisioner{}
}

const firectlApiKeyFlag = "--api-key"

func (p firectlCLIArgProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	apiKey, ok := in.ItemFields[fieldname.APIKey]
	if !ok {
		out.AddError(errors.New("missing API key field"))
		return
	}

	out.AddArgs(firectlApiKeyFlag, apiKey)
}

func (p firectlCLIArgProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// No op
}

func (p firectlCLIArgProvisioner) Description() string {
	return "Firectl CLI API Key argument provisioner"
}
