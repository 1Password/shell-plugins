package provisioner

import (
	"context"
	"github.com/1Password/shell-plugins/plugins/age/operation"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

// KeyFiles holds the private and public key material as file contents.
type KeyFiles struct {
	private provision.ItemToFileContents
	public  provision.ItemToFileContents
}

// NewKeyFiles creates a new instance of KeyFiles to hold the key material.
func NewKeyFiles(private, public provision.ItemToFileContents) *KeyFiles {
	return &KeyFiles{private, public}
}

// KeyPairProvisioner handles the provisioning and deprovisioning of key files for the age command.
type KeyPairProvisioner struct {
	private PrivateKeyProvisioner
	public  PublicKeyProvisioner
}

// Provision sets up the necessary key file(s) for the `age` command based on the operation mode (i.e. encrypt or decrypt).
// It determines the mode from the command line arguments and calls the appropriate provisioner.
func (p KeyPairProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	switch operation.Detect(out.CommandLine) {
	case operation.Encrypt:
		p.public.Provision(ctx, in, out)
	case operation.Decrypt:
		p.private.Provision(ctx, in, out)

	}
}

// Deprovision performs cleanup after the process completes.
// In this implementation, no cleanup is required as temporary files are automatically removed.
func (p KeyPairProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

func (p KeyPairProvisioner) Description() string {
	return "Determine operating mode and provision appropriate key files and populate command line arguments."
}

// KeyPairTempFile creates a new KeyPairProvisioner for handling temporary files for the age command.
func KeyPairTempFile(keys *KeyFiles) sdk.Provisioner {
	return KeyPairProvisioner{
		private: PrivateKeyTempFile(keys.private),
		public:  PublicKeyTempFile(keys.public),
	}
}
