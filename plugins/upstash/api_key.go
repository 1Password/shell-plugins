package upstash

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
		DocsURL:       sdk.URL("https://docs.upstash.com/redis/account/developerapi#create-an-api-key"),
		ManagementURL: sdk.URL("https://console.upstash.com/account/api"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Upstash.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
			{
				Name:                fieldname.Email,
				MarkdownDescription: "Email used to authenticate to Upstash.",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryUpstashConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"UPSTASH_API_KEY": fieldname.APIKey,
	"UPSTASH_EMAIL":   fieldname.Email,
}

func TryUpstashConfigFile() sdk.Importer {
	return importer.TryFile("~/.upstash.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" && config.Email == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: config.APIKey,
				fieldname.Email:  config.Email,
			},
		})
	})
}

type Config struct {
	APIKey string `json:"apiKey"`
	Email  string `json:"email"`
}
