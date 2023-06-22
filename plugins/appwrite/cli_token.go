package appwrite

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func CLIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.CLIToken,
		DocsURL:       sdk.URL("https://appwrite.io/docs/command-line"),
		ManagementURL: sdk.URL("https://cloud.appwrite.io/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Email,
				MarkdownDescription: "Email used to authenticate to Appwrite CLI.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 50,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Password,
				MarkdownDescription: "Password used to authenticate to Appwrite CLI.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 23,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
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
	"APPWRITE_EMAIL":    fieldname.Email,
	"APPWRITE_PASSWORD": fieldname.Password,
}
