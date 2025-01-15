package age

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/provision"
)

// KeyFiles holds the private and public key material as file contents.
type KeyFiles struct {
	private provision.ItemToFileContents
	public  provision.ItemToFileContents
}

// KeyPairProvisioner handles the provisioning and deprovisioning of key files for the age command.
type KeyPairProvisioner struct {
	keys        KeyFiles
	fileOptions []provision.FileOption
}

// Provision sets up the necessary key file for the `age` command based on the operation mode (Encrypt or Decrypt).
// It determines the mode from the command line arguments and prepares the corresponding key file.
func (p KeyPairProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	operationHandlers := map[Operation]OperationHandler{
		Encrypt: handleEncrypt,
		Decrypt: handleDecrypt,
	}

	mode := detectOperation(out.CommandLine)
	handler, ok := operationHandlers[mode]
	if !ok {
		out.AddError(ErrUnknownCommand)
		return
	}

	keyFileMaterialiser, args, filename := handler(p.keys, out)

	p.fileOptions = append(p.fileOptions, provision.Filename(filename), provision.PrependArgs(args...))

	fileProvisioner := provision.TempFile(keyFileMaterialiser, p.fileOptions...)
	fileProvisioner.Provision(ctx, in, out)
}

// Deprovision performs cleanup after the process completes.
// In this implementation, no cleanup is required as temporary files are automatically removed.
func (p KeyPairProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

func (p KeyPairProvisioner) Description() string {
	return "Provision temporary file with public & private key pair & pass to age command"
}

func handleEncrypt(keys KeyFiles, _ *sdk.ProvisionOutput) (provision.ItemToFileContents, []string, string) {
	return keys.public, []string{"-R", "{{.Path}}"}, "age.public.txt"
}

func handleDecrypt(keys KeyFiles, out *sdk.ProvisionOutput) (provision.ItemToFileContents, []string, string) {
	for _, arg := range out.CommandLine {
		if arg == "-i" || arg == "--identity" {
			out.AddError(ErrConflictingIdentityFlag)
		}
	}
	return keys.private, []string{"-i", "{{.Path}}"}, "age.private.txt"
}

// TempFile creates a new KeyPairProvisioner for handling temporary files for the age command.
func TempFile(keys KeyFiles, opts ...provision.FileOption) sdk.Provisioner {
	return KeyPairProvisioner{
		keys:        keys,
		fileOptions: opts,
	}
}
