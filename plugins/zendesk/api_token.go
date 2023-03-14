package zendesk

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.APIToken,
		DocsURL: sdk.URL("https://developer.zendesk.com/api-reference/introduction/security-and-auth/#api-token"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.OrgURL,
				MarkdownDescription: "Subdomain of Zendesk account, often found in the account's URL.",
				Optional:            false,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Email used to authenticate to Zendesk.",
				Optional:            false,
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API token used to authenticate to Zendesk.",
				Optional:            false,
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
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
	"ZENDESK_SUBDOMAIN": fieldname.OrgURL,
	"ZENDESK_EMAIL":     fieldname.Username,
	"ZENDESK_API_TOKEN": fieldname.Token,
}
