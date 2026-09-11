package provisioner

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

// PrivateKeyProvisioner handles the provisioning and deprovisioning of private key files for the age command.
type PrivateKeyProvisioner struct {
	privateKey  provision.ItemToFileContents
	fileOptions []provision.FileOption
}

// Provision sets up the private key file for age decrypt commands.
func (p PrivateKeyProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	for _, arg := range out.CommandLine {
		if arg == "-i" || arg == "--identity" {
			out.AddError(ErrConflictingIdentityFlag)
		}
	}

	fileProvisioner := provision.TempFile(p.privateKey, p.fileOptions...)
	fileProvisioner.Provision(ctx, in, out)
}

// Deprovision performs cleanup after the process completes.
// In this implementation, no cleanup is required as temporary files are automatically removed.
func (p PrivateKeyProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

func (p PrivateKeyProvisioner) Description() string {
	return "Provision temporary file with private key & populate command line arguments."
}

// PrivateKeyTempFile creates a new PrivateKeyProvisioner for creating private key files for age.
func PrivateKeyTempFile(privateKey provision.ItemToFileContents, opts ...provision.FileOption) PrivateKeyProvisioner {
	opts = append(opts, provision.Filename("age.private.txt"), provision.PrependArgs("-i", "{{.Path}}"))
	return PrivateKeyProvisioner{
		privateKey:  privateKey,
		fileOptions: opts,
	}
}
