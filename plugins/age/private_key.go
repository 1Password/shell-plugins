package age

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PrivateKey() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.SecretKey,
		DocsURL: sdk.URL("https://age-encryption.org/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Age Private key",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 80,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: PrivateKeyProvisioner{},
	}
}
