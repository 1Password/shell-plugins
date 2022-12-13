package datadog

import (
	"context"

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

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DATADOG_API_KEY": fieldname.APIKey,
	"DATADOG_APP_KEY": fieldname.AppKey,
}

func TryDogrcFile() sdk.Importer {
	return importer.TryFile("~/.dogrc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("apikey") && section.Key("apikey").Value() != "" {
				fields[fieldname.APIKey] = section.Key("apikey").Value()
			}

			if section.HasKey("appkey") && section.Key("appkey").Value() != "" {
				fields[fieldname.AppKey] = section.Key("appkey").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}
