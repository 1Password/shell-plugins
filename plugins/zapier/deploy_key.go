package zapier

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func DeployKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.DeployKey,
		DocsURL:       sdk.URL("https://platform.zapier.com/cli_docs/docs#quick-setup-guide"),
		ManagementURL: sdk.URL("https://developer.zapier.com/partner-settings/deploy-keys/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Key,
				MarkdownDescription: "Deploy Key used to authenticate to Zapier.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
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
			TryZapierConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"ZAPIER_DEPLOY_KEY": fieldname.Key,
}

func TryZapierConfigFile() sdk.Importer {
	return importer.TryFile("~/.zapierrc", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.DeployKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Key: config.DeployKey,
			},
		})
	})
}

type Config struct {
	DeployKey string `json:"deployKey"`
}
