package wrangler

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
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://developers.cloudflare.com/fundamentals/api/get-started/create-token/"),
		ManagementURL: sdk.URL("https://dash.cloudflare.com/profile/api-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AccountID,
				MarkdownDescription: "The account ID for the Workers related account, can be found in the Cloudflare dashboard, can usually be inferred by Wrangler.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "The API token for your Cloudflare account, can be used for authentication for situations like CI/CD, and other automation.",
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
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer:           importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CLOUDFLARE_ACCOUNT_ID": fieldname.AccountID,
	"CLOUDFLARE_API_TOKEN":  fieldname.Token,
}
