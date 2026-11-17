package octopus

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
		DocsURL:       sdk.URL("https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key"),
		ManagementURL: nil,
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Octopus Deploy.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Digits:    true,
					},
					Prefix: "API-",
				},
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "URL of the Octopus Deploy server.",
				Optional:            true,
			},
			{
				Name:                fieldname.Space,
				MarkdownDescription: "Space to use for the Octopus Deploy CLI.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOctopusDeployConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OCTOPUS_API_KEY": fieldname.APIKey,
	"OCTOPUS_URL":     fieldname.URL,
	"OCTOPUS_SPACE":   fieldname.Space,
}

func TryOctopusDeployConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/octopus/cli_config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" {
			return
		}

		fields := map[sdk.FieldName]string{
			fieldname.APIKey: config.APIKey,
		}

		if config.URL != "" {
			fields[fieldname.URL] = config.URL
		}

		if config.Space != "" {
			fields[fieldname.Space] = config.Space
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}

type Config struct {
	APIKey string `json:"apikey"`
	URL    string `json:"url"`
	Space  string `json:"space"`
}
