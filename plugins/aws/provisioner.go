package aws

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type awsProvisioner struct {
	stsProvisioner    STSProvisioner
	envVarProvisioner provision.EnvVarProvisioner
}

func AWSProvisioner(envVarToFieldName map[string]string) sdk.Provisioner {
	return awsProvisioner{
		envVarProvisioner: provision.EnvVarProvisioner{
			Schema: envVarToFieldName,
		},
	}
}

func (p awsProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	totpCode, foundTotp := in.ItemFields[fieldname.OneTimePassword]
	serialNumber, foundSerialNumber := in.ItemFields[FieldNameSerialNumber]
	if foundTotp && foundSerialNumber {
		p.stsProvisioner.TOTPCode = totpCode
		p.stsProvisioner.MFASerial = serialNumber
		p.stsProvisioner.Provision(ctx, in, out)
	} else {
		p.envVarProvisioner.Provision(ctx, in, out)
	}
}

func (p awsProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p awsProvisioner) Description() string {
	return fmt.Sprintf("%s, and, if MFA is present, %s", p.envVarProvisioner.Description(), p.stsProvisioner.Description())
}
