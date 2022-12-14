package circleci

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAPIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAPIToken,
		DocsURL:       sdk.URL("https://circleci.com/docs/managing-api-tokens"),
		ManagementURL: sdk.URL("https://app.circleci.com/settings/user/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to CircleCI.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(map[string]sdk.FieldName{
			"CIRCLECI_CLI_TOKEN": fieldname.Token,
		}),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"CIRCLECI_CLI_TOKEN": fieldname.Token,
			}),
			TryCircleCIConfigFile(),
		),
	}
}

type Config struct {
	Token string `yaml:"token"`
}

func TryCircleCIConfigFile() sdk.Importer {
	return importer.TryFile("~/.circleci/cli.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.Token == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: config.Token,
			},
		})
	})
}
