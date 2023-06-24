package civo

import (
	"context"
	"encoding/json"

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
			{
				Name:                fieldname.APIKeyID,
				MarkdownDescription: "API Name to identify the API Key.",
			},
			{
				Name:                fieldname.DefaultRegion,
				MarkdownDescription: "The default region to use for this API Key.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryCivoConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"CIVO_TOKEN":        fieldname.APIKey,
	"CIVO_API_KEY_NAME": fieldname.APIKeyID,
	"CIVO_API_KEY":      fieldname.APIKey,
}

func TryCivoConfigFile() sdk.Importer {

	return importer.TryFile("~/.civo.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return

		}

		if len(config.Properties) == 0 && config.Meta.CurrentAPIKey == "" {
			return
		}

		for key, value := range config.Properties {
			var apiKey string

			err := json.Unmarshal(value, &apiKey)
			if err != nil {
				out.AddError(err)
				return
			}
			out.AddCandidate(sdk.ImportCandidate{
				NameHint: key,
				Fields: map[sdk.FieldName]string{
					fieldname.APIKey:        apiKey,
					fieldname.APIKeyID:      config.Meta.CurrentAPIKey,
					fieldname.DefaultRegion: config.Meta.DefaultRegion,
				},
			})

			break
		}

	})
}

type Config struct {
	Properties map[string]json.RawMessage `json:"apikeys"`

	Meta struct {
		CurrentAPIKey string `json:"current_apikey"`
		DefaultRegion string `json:"default_region"`
	} `json:"meta"`
}
