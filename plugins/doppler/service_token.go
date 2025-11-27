package doppler

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func ServiceToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.doppler.com/docs/cli"),
		ManagementURL: sdk.URL("https://dashboard.doppler.com/workplace/<workplace-id>/tokens/personal"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Personal Token used to authenticate to Doppler.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Prefix: "dp.pt.",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryAllEnvVars(fieldname.Token, "DOPPLER_TOKEN"),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DOPPLER_TOKEN": fieldname.Token,
}
