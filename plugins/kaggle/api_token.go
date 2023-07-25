package kaggle

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://www.kaggle.com/docs/api"),
		ManagementURL: sdk.URL("https://www.kaggle.com/settings/account"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API Token used to authenticate to Kaggle.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Username,
				MarkdownDescription: "Username to authenticate to Kaggle.",
				Secret:              false,
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
			TryKaggleConfigFile("~/.kaggle/kaggle.json"),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"KAGGLE_KEY":      fieldname.Token,
	"KAGGLE_USERNAME": fieldname.Username,
}

func TryKaggleConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.Token == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token:    config.Token,
				fieldname.Username: config.Username,
			},
			NameHint: importer.SanitizeNameHint(config.Username),
		})
	})
}

type Config struct {
	Username string `json:"username"`
	Token    string `json:"key"`
}
