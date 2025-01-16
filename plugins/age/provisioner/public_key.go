package provisioner

import (
	"context"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

// PublicKeyProvisioner handles the provisioning and deprovisioning of public key files for the age command.
type PublicKeyProvisioner struct {
	publicKey   provision.ItemToFileContents
	fileOptions []provision.FileOption
}

// Provision sets up the public key file for age encrypt commands.
func (p PublicKeyProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	fileProvisioner := provision.TempFile(p.publicKey, p.fileOptions...)
	fileProvisioner.Provision(ctx, in, out)
}

// Deprovision performs cleanup after the process completes.
// In this implementation, no cleanup is required as temporary files are automatically removed.
func (p PublicKeyProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

func (p PublicKeyProvisioner) Description() string {
	return "Provision temporary file with public key file & populate command line arguments."
}

// PublicKeyTempFile creates a new PublicKeyProvisioner for creating public key files for age.
func PublicKeyTempFile(publicKey provision.ItemToFileContents, opts ...provision.FileOption) PublicKeyProvisioner {
	opts = append(opts, provision.Filename("age.public.txt"), provision.PrependArgs("-R", "{{.Path}}"))
	return PublicKeyProvisioner{
		publicKey:   publicKey,
		fileOptions: opts,
	}
}
