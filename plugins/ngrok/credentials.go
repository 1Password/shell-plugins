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

func Credentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.Credentials,
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/config"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AuthToken,
				MarkdownDescription: "Auth Token used to authenticate to ngrok.",
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
		DefaultProvisioner: provision.TempFile(ngrokConfig, provision.Filename("ngrok.yml"), provision.AddArgs("--config", "{{ .Path }}")),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryngrokConfigFile(),
		)}
}

func ngrokConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		AuthToken: in.ItemFields[fieldname.AuthToken],
		APIKey:    in.ItemFields[fieldname.APIKey],
		Version:   "2", // required field for ngrok CLI to work
	}
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NGROK_AUTHTOKEN": fieldname.AuthToken,
	"NGROK_API_KEY":   fieldname.APIKey,
}

func TryngrokConfigFile() sdk.Importer {
	return importer.TryFile("~/.ngrok2/ngrok.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AuthToken == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.AuthToken: config.AuthToken,
				fieldname.APIKey:    config.APIKey,
			},
		})
	})
}

type Config struct {
	AuthToken string `yaml:"authtoken"`
	APIKey    string `yaml:"api_key"`
	Version   string
}
