package age

import (
	"context"
	"crypto/rand"
	"fmt"
	"slices"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

type PrivateKeyProvisioner struct {
}

func (p PrivateKeyProvisioner) Provision(_ context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	privateKey := in.ItemFields[fieldname.PrivateKey]

	fileName, err := randomFilename()
	if err != nil {
		out.AddError(fmt.Errorf("generating random file name: %s", err))
		return
	}

	secretFilePath := in.FromTempDir(fileName)
	out.AddSecretFile(secretFilePath, []byte(privateKey))

	editedCommandLine := slices.Insert(out.CommandLine, 1, "--identity", secretFilePath)

	if editedCommandLine != nil {
		out.CommandLine = editedCommandLine
	}
}

func (p PrivateKeyProvisioner) Deprovision(_ context.Context, _ sdk.DeprovisionInput, _ *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (p PrivateKeyProvisioner) Description() string {
	return "Provision age secret key"
}

func randomFilename() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
