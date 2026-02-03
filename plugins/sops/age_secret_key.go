package sops

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AgeSecretKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.SecretKey,
		DocsURL:       sdk.URL("https://github.com/getsops/sops#encrypting-using-age"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.PrivateKey,
				MarkdownDescription: "Age secret key used by SOPS for encryption and decryption.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Prefix: "AGE-SECRET-KEY-",
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SOPS_AGE_KEY": fieldname.PrivateKey,
}
