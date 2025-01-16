package age

import (
	"fmt"

	"github.com/1Password/shell-plugins/plugins/age/provisioner"
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func KeyPair() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.SecretKey,
		DocsURL: sdk.URL("https://age-encryption.org/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.PublicKey,
				MarkdownDescription: "Age X25519 public key.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 62,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Age X25519 private key.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 74,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
		},
		DefaultProvisioner: provisioner.KeyPairTempFile(provisioner.NewKeyFiles(
			materialisePrivateKeyFile,
			materialisePublicKeyFile,
		)),
		Importer: importer.NoOp(),
	}
}

func materialisePublicKeyFile(in sdk.ProvisionInput) ([]byte, error) {
	content := "# generated by: 1password-cli/shell-plugins/plugin/age\n"

	if publicKey, ok := in.ItemFields[fieldname.PublicKey]; ok {
		content += publicKey
	}

	return []byte(content), nil
}

func materialisePrivateKeyFile(in sdk.ProvisionInput) ([]byte, error) {
	content := "# generated by: 1password-cli/shell-plugins/plugin/age\n"

	if publicKey, ok := in.ItemFields[fieldname.PublicKey]; ok {
		content += fmt.Sprintf("# public key: %s\n", publicKey)
	}

	if privateKey, ok := in.ItemFields[fieldname.PrivateKey]; ok {
		content += privateKey
	}

	return []byte(content), nil
}
