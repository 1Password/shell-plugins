package civo

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
		DocsURL:       sdk.URL("https://www.civo.com/docs/account/api-keys"),
		ManagementURL: sdk.URL("https://dashboard.civo.com/security"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Civo.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 50,
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
			TryCivoConfigFile("~/.civo.json"),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CIVO_TOKEN": fieldname.APIKey,
}

func TryCivoConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return

		}

		if len(config.Properties) == 0 {
			return
		}

		for key, value := range config.Properties {
			out.AddCandidate(sdk.ImportCandidate{
				NameHint: key,
				Fields: map[sdk.FieldName]string{
					fieldname.APIKey: value,
				},
			})
		}
	})
}

type Config struct {
	Properties map[string]string `json:"apikeys"`
}
