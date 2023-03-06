package ngrok

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v3"
)

func APICredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.APIKey,
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/config"),
		Fields: []schema.CredentialField{
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
		DefaultProvisioner: provision.TempFile(ngrokApiConfig, provision.Filename("config.yml"), provision.AddArgs("--config", "{{ .Path }}")),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultApiEnvVarMapping),
			importer.MacOnly(
				TryngrokAPIConfigFile("~/Library/Application Support/ngrok/ngrok.yml"),
			),
			importer.LinuxOnly(
				TryngrokAPIConfigFile("~/.config/ngrok/ngrok.yml"),
			),
		)}
}

func ngrokApiConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := APIConfig{
		APIKey:  in.ItemFields[fieldname.APIKey],
		Version: "2", // required field for ngrok CLI to work when file-based configuration is used; automatically configured by the CLI program and is not configurable by the user
	}
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

var defaultApiEnvVarMapping = map[string]sdk.FieldName{
	"NGROK_API_KEY": fieldname.APIKey,
}

func TryngrokAPIConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config APIConfig
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" || config.Version == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: config.APIKey,
			},
		})
	})
}

type APIConfig struct {
	APIKey  string `yaml:"api_key"`
	Version string
}
