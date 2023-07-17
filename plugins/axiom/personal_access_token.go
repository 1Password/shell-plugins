package axiom

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://axiom.co/docs/restapi/token#creating-personal-token"),
		ManagementURL: sdk.URL("https://app.axiom.co/settings/profile"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Axiom.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 41,
					Prefix: "xapt-",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "The organization ID of the organization the access token is valid for. Only valid for Axiom Cloud.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
				},
			},
			{
				Name:                fieldname.Deployment,
				MarkdownDescription: "Deployment to use.",
				Secret:              false,
				Optional:            true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
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
	"AXIOM_TOKEN":      fieldname.Token,
	"AXIOM_ORG_ID":     fieldname.Organization,
	"AXIOM_DEPLOYMENT": fieldname.Deployment,
}
