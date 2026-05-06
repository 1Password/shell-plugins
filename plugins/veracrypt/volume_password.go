package veracrypt

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"

	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func VolumePassword() schema.CredentialType {
	return schema.CredentialType{
		Name:          sdk.CredentialName("Volume Password"),
		DocsURL:       sdk.URL("https://www.veracrypt.fr/en/Documentation.html"),
		ManagementURL: sdk.URL("https://www.veracrypt.fr/en/Main.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to mount a VeraCrypt volume.",
				Secret:              true,
			},
			{
				Name:                sdk.FieldName("Volume"),
				MarkdownDescription: "Path to the VeraCrypt volume file.",
				Secret:              false,
				Optional:            true,
			},
		},
		DefaultProvisioner: volumePasswordProvisioner(),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryVeraCryptConfigFile(),
		),
	}
}

type volumePasswordProv struct{}

func volumePasswordProvisioner() sdk.Provisioner {
	return volumePasswordProv{}
}

func (p volumePasswordProv) Description() string {
	return "Provision password as command-line arguments"
}

func (p volumePasswordProv) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	password, ok := in.ItemFields[fieldname.Password]
	if !ok || password == "" {
		out.AddError(fmt.Errorf("password is required"))
		out.CommandLine = []string{}
		return
	}
	out.AddArgs("-p", password, "--non-interactive")
}

func (p volumePasswordProv) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"VERACRYPT_PASSWORD": fieldname.Password,
}

func TryVeraCryptConfigFile() sdk.Importer {
	return importer.TryFile("~/.VeraCrypt/Config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
	})
}