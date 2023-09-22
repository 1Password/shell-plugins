package ansiblevault

import (
	"context"
	"fmt"
	"os"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Password() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.Credential,
		DocsURL:       sdk.URL("https://docs.ansible.com/ansible/latest/cli/ansible-vault.html"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password to use when encrypting or decrypting.",
				Secret:              true,
			},
		},
		DefaultProvisioner: provision.TempFile(
			passwordFile,
			provision.Filename(".ansible-vault"),
			provision.SetPathAsEnvVar("ANSIBLE_VAULT_PASSWORD_FILE"),
		),
		Importer: importer.TryAll(
			TryPasswordFile(),
		)}
}

func passwordFile(in sdk.ProvisionInput) ([]byte, error) {
	if password, ok := in.ItemFields[fieldname.Password]; ok {
		return []byte(password + "\n"), nil
	}
	return nil, fmt.Errorf("unable to find password field")
}

func TryPasswordFile() sdk.Importer {
	file := os.Getenv("ANSIBLE_VAULT_PASSWORD_FILE")
	if file == "" {
		file = "~/.ansible-vault"
	}
	return importer.TryFile(file, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		password := contents.ToString()
		if password == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Password: password,
			},
		})
	})
}
