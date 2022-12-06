package twilio

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
		DocsURL:       sdk.URL("https://www.twilio.com/docs/glossary/what-is-an-api-key"),
		ManagementURL: sdk.URL("https://www.twilio.com/console/project/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccountSID,
				MarkdownDescription: "Account SID used to authenticate to Twilio.",
				Composition: &schema.ValueComposition{
					Length: 34,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
					Prefix: "AC",
				},
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Twilio.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 34,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
					Prefix: "SK",
				},
			},
			{
				Name:                fieldname.APISecret,
				MarkdownDescription: "API Secret used to authenticate to Twilio.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Region,
				MarkdownDescription: "The region to use for this API Key.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"TWILIO_ACCOUNT_SID": fieldname.AccountSID,
	"TWILIO_API_KEY":     fieldname.APIKey,
	"TWILIO_API_SECRET":  fieldname.APISecret,
}
