package prefect

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.prefect.io/2.10.13/cloud/users/api-keys/"),
		ManagementURL: sdk.URL("https://app.prefect.cloud/my/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "Key used to authenticate to Prefect.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "pnu_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "URL for the Prefect workspace.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Length: 124,
					Prefix: "https://api.prefect.cloud/api/accounts/",
					Charset: schema.Charset{
						Uppercase: false,
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
	"PREFECT_API_KEY": fieldname.APIKey,
	"PREFECT_API_URL": fieldname.URL,
}
