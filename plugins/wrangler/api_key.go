package wrangler

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
		DocsURL:       sdk.URL("https://developers.cloudflare.com/workers/wrangler/system-environment-variables/"),
		ManagementURL: sdk.URL("https://dash.cloudflare.com/profile/api-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Username,
				MarkdownDescription: "The email address associated with your Cloudflare account, usually used for older authentication method with CLOUDFLARE_API_KEY",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "The API key for your Cloudflare account, usually used for older authentication method with CLOUDFLARE_EMAIL",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMappingApiKey),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMappingApiKey),
	}
}

var defaultEnvVarMappingApiKey = map[string]sdk.FieldName{
	"CLOUDFLARE_EMAIL":   fieldname.Username,
	"CLOUDFLARE_API_KEY": fieldname.APIKey,
}
