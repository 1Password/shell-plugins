package exercism

import (
	"context"
	"encoding/json"

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
		DocsURL:       sdk.URL("https://exercism.org/docs/using/solving-exercises/working-locally"),
		ManagementURL: sdk.URL("https://exercism.org/settings/api_cli"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.URL,
				MarkdownDescription: `The URL of the Exercism API.`,
				Secret:              false,
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Exercism.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 37,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                sdk.FieldName("Directory"),
				MarkdownDescription: `The path to the workspace directory where the exercises are stored.`,
				Secret:              false,
			},
		},
		DefaultProvisioner: provision.TempFile(
			tempFileConfig,
			provision.AtFixedPath("~/.config/exercism/user.json"),
		),
		Importer: importer.TryAll(
			TryExercismConfigFile(),
		)}
}

func TryExercismConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/exercism/user.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.URL == "" || config.APIKey == "" || config.Workspace == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey:           config.APIKey,
				fieldname.URL:              config.URL,
				sdk.FieldName("Directory"): config.Workspace,
			},
		})
	})
}

type Config struct {
	URL       string `json:"apibaseurl"`
	APIKey    string `json:"token"`
	Workspace string `json:"workspace"`
}

func tempFileConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		URL:       in.ItemFields[fieldname.URL],
		APIKey:    in.ItemFields[fieldname.APIKey],
		Workspace: in.ItemFields[sdk.FieldName("Directory")],
	}

	contents, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}
