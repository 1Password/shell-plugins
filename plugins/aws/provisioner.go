package aws

import (
	"context"
	"fmt"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type AWSProvisioner struct {
	stsProvisioner    StsProvisioner
	envVarProvisioner provision.EnvVarProvisioner
}

func getProvisioner(envVarToFieldName map[string]string) sdk.Provisioner {
	return AWSProvisioner{
		envVarProvisioner: provision.EnvVarProvisioner{
			Schema: envVarToFieldName,
		},
	}
}

func (p AWSProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	totpCode, foundTotp := in.ItemFields[fieldname.OneTimePassword]
	serialNumber, foundSerialNumber := in.ItemFields[FieldNameSerialNumber]
	if foundTotp && foundSerialNumber {
		p.stsProvisioner.TotpCode = totpCode
		p.stsProvisioner.SerialNumber = serialNumber
		p.stsProvisioner.Provision(ctx, in, out)
	} else {
		p.envVarProvisioner.Provision(ctx, in, out)
	}
}

func (p AWSProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p AWSProvisioner) Description() string {
	return fmt.Sprintf("%s, and, if MFA is present, %s", p.envVarProvisioner.Description(), p.stsProvisioner.Description())
}
