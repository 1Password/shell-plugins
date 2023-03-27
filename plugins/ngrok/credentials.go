package ngrok

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.Credentials,
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/config"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Authtoken,
				AlternativeNames:    []string{"Auth Token"},
				MarkdownDescription: "Authtoken used to authenticate to ngrok.",
				Optional:            false,
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 43,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to ngrok API.",
				Optional:            true,
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 48,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: ngrokProvisioner(),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.MacOnly(
				TryngrokConfigFile("~/Library/Application Support/ngrok/ngrok.yml"),
			),
			importer.LinuxOnly(
				TryngrokConfigFile("~/.config/ngrok/ngrok.yml"),
			),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NGROK_AUTHTOKEN": fieldname.Authtoken,
	"NGROK_API_KEY":   fieldname.APIKey,
}

func TryngrokConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AuthToken == "" || config.APIKey == "" || config.Version == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Authtoken: config.AuthToken,
				fieldname.APIKey:    config.APIKey,
			},
		})
	})
}

// Config struct is exhaustive, covering all documented configurations.
type Config struct {
	AuthToken string `yaml:"authtoken"`
	APIKey    string `yaml:"api_key"`
	Version   string `yaml:"version"`
}
