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

func AuthCredentials() schema.CredentialType {
	return schema.CredentialType{
		Name:    credname.AuthToken,
		DocsURL: sdk.URL("https://ngrok.com/docs/ngrok-agent/config"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.AuthToken,
				AlternativeNames:    []string{"Auth Token"},
				MarkdownDescription: "AuthToken used to authenticate to ngrok.",
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
		},
		DefaultCredential:  true,
		DefaultProvisioner: provision.TempFile(ngrokConfig, provision.Filename("config.yml"), provision.AddArgs("--config", "{{ .Path }}")),
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

func ngrokConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		AuthToken: in.ItemFields[fieldname.AuthToken],
		Version:   "2", // required field for ngrok CLI to work when file-based configuration is used; automatically configured by the CLI program and is not configurable by the user
	}
	contents, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NGROK_AUTHTOKEN": fieldname.AuthToken,
}

func TryngrokConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AuthToken == "" || config.Version == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.AuthToken: config.AuthToken,
			},
		})
	})
}

type Config struct {
	AuthToken string `yaml:"authtoken"`
	Version   string
}
