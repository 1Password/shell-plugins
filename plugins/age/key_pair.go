package age

import (
	"fmt"
	"github.com/1Password/shell-plugins/sdk/provision"

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
		DefaultProvisioner: provision.NoOp(),
		Importer:           importer.NoOp(),
	}
}
