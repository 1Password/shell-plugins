package fossa

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
		DocsURL:       sdk.URL("https://docs.fossa.com/docs/api-reference#api-tokens"),
		ManagementURL: sdk.URL("https://app.fossa.com/account/settings/integrations/api_tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to FOSSA.",
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
		Importer:    importer.TryEnvVarPair(defaultEnvVarMapping),
	}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.APIKey: "FOSSA_API_KEY",
}
