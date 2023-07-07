package b2

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ApplicationKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.ApplicationKey,
		DocsURL:       sdk.URL("https://www.backblaze.com/docs/cloud-storage-application-keys"),
		ManagementURL: sdk.URL("https://secure.backblaze.com/app_keys.htm"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.ApplicationKeyID,
				MarkdownDescription: "The key ID of the application key used to authenticate to Backblaze B2.",
				Composition: &schema.ValueComposition{
					Length: 25,
					Prefix: "003",
					Charset: schema.Charset{
						Specific: []rune("0123456789abcdef"),
					},
				},
			},
			{
				Name:                fieldname.ApplicationKey,
				MarkdownDescription: "The secret part of the application key used to authenticate to Backblaze B2.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 25,
					Prefix: "K003",
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Specific:  []rune{'/'},
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
	"B2_ACCOUNT_ID":  fieldname.ApplicationKeyID, // TODO: Check if this is correct
	"B2_ACCOUNT_KEY": fieldname.ApplicationKey,   // TODO: Check if this is correct
}
