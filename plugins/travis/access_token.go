package travis

import (
	"context"
	"net/url"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://developer.travis-ci.com/authentication"),
		ManagementURL: sdk.URL("https://app.travis-ci.com/account/preferences"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Travis CI.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 22,
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
			TryTravisCIConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"TRAVIS_TOKEN": fieldname.Token,
}

func TryTravisCIConfigFile() sdk.Importer {
	return importer.TryFile("~/.travis/config.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for endpoint, endpointConfig := range config.Endpoints {
			if endpointConfig.AccessToken == "" {
				continue
			}

			nameHint := endpoint
			if parsed, err := url.Parse(endpoint); err == nil && parsed.Host != "" {
				nameHint = parsed.Host
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: endpointConfig.AccessToken,
				},
				NameHint: importer.SanitizeNameHint(nameHint),
			})
		}
	})
}

type Config struct {
	Endpoints map[string]Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	AccessToken string `yaml:"access_token"`
}
