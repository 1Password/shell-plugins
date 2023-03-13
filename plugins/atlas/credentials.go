package atlas

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name: credname.Credentials,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.PublicKey,
				MarkdownDescription: "Public key used to authenticate to MongoDB Atlas.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 8,
					Charset: schema.Charset{
						Lowercase: true,
					},
				},
			},
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Private key used to authenticate to MongoDB Atlas.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"MONGODB_ATLAS_PUBLIC_API_KEY":  fieldname.PublicKey,
	"MONGODB_ATLAS_PRIVATE_API_KEY": fieldname.PrivateKey,
}
