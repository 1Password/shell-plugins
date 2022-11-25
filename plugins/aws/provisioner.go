package aws

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type awsProvisioner struct {
	stsProvisioner    STSProvisioner
	envVarProvisioner provision.EnvVarProvisioner
}

func AWSProvisioner() sdk.Provisioner {
	return awsProvisioner{
		envVarProvisioner: provision.EnvVarProvisioner{
			Schema: officialEnvVarMapping,
		},
	}
}

func (p awsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	totp, hasTotp := in.ItemFields[fieldname.OneTimePassword]
	mfaSerial, hasMFASerial := in.ItemFields[FieldNameSerialNumber]

	if hasTotp && hasMFASerial {
		p.stsProvisioner.MFASerial = mfaSerial
		p.stsProvisioner.TOTPCode = totp
		p.stsProvisioner.Provision(ctx, in, out)
	} else {
		p.envVarProvisioner.Provision(ctx, in, out)
	}
}

func (p awsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p awsProvisioner) Description() string {
	return p.envVarProvisioner.Description()
}
