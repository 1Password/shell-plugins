package spacelift

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
		DocsURL:       sdk.URL("https://docs.spacelift.io/integrations/api#spacelift-api-key-token"),
		ManagementURL: sdk.URL("https://mycorp.app.spacelift.io/settings/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Endpoint,
				MarkdownDescription: "The URL to your Spacelift account, for example https://mycorp.app.spacelift.io",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Symbols:   true,
					},
					Prefix: "https://",
				},
			},
			{
				Name:                fieldname.APIKeyID,
				MarkdownDescription: "The ID of your Spacelift API key. Available via the Spacelift application.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.APIKeySecret,
				MarkdownDescription: "The secret for your API key. Only available when the secret is created.",
				Secret:              true,
				Composition: &schema.ValueComposition{
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
			TrySpaceliftConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SPACELIFT_API_KEY_ENDPOINT": fieldname.Endpoint,
	"SPACELIFT_API_KEY_ID":       fieldname.APIKeyID,
	"SPACELIFT_API_KEY_SECRET":   fieldname.APIKeySecret,
}

func TrySpaceliftConfigFile() sdk.Importer {
	return importer.TryFile("~/.spacelift/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		for name, profile := range config.Profiles {
			if profile.Credentials.Type != 1 {
				continue
			}

			out.AddCandidate(sdk.ImportCandidate{
				NameHint: name,
				Fields: map[sdk.FieldName]string{
					fieldname.Endpoint:     profile.Credentials.Endpoint,
					fieldname.APIKeyID:     profile.Credentials.KeyID,
					fieldname.APIKeySecret: profile.Credentials.KeySecret,
				},
			})
		}
	})
}

type Config struct {
	Profiles map[string]SpaceCLIProfile `json:"profiles"`
}

type SpaceCLIProfile struct {
	Credentials SpaceCLIProfileCredentials `json:"credentials"`
}

type SpaceCLIProfileCredentials struct {
	Endpoint  string `json:"endpoint"`
	Type      int    `json:"type"`
	KeyID     string `json:"key_id"`
	KeySecret string `json:"key_secret"`
}
