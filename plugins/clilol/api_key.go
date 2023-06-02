package clilol

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		ManagementURL: sdk.URL("https://home.omg.lol/account#api-key"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Address,
				MarkdownDescription: "Your omg.lol address.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 63,
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Email,
				MarkdownDescription: "Your omg.lol account email address.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Specific:  []rune{'@', '.'},
					},
				},
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to omg.lol.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
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
	"CLILOL_ADDRESS": fieldname.Address,
	"CLILOL_EMAIL":   fieldname.Email,
	"CLILOL_APIKEY":  fieldname.APIKey,
}
