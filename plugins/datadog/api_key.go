package datadog

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
		DocsURL:       sdk.URL("https://docs.datadoghq.com/account_management/api-app-keys/"),
		ManagementURL: sdk.URL("https://app.datadoghq.com/organization-settings/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Datadog.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.AppKey,
				MarkdownDescription: "Application key used for this API key.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
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
			TryDogrcFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.APIKey: "DATADOG_API_KEY",
	fieldname.AppKey: "DATADOG_APP_KEY",
}

func TryDogrcFile() sdk.Importer {
	// TODO: Try importing from ~/.dogrc file
	return importer.NoOp()
}
