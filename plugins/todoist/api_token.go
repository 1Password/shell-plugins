package todoist

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

func APIToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIToken,
		DocsURL:       sdk.URL("https://todoist.com/help/articles/8048880904476"),
		ManagementURL: sdk.URL("https://todoist.com/app/settings/integrations/developer"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API Token used to authenticate to Todoist.",
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
		DefaultProvisioner: provision.TempFile(
			todoistConfig,
			provision.AtFixedPath("~/.config/todoist/config.json"),
		),
		Importer: importer.TryAll(
			TryTodoistConfigFile(),
		)}
}

func TryTodoistConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/todoist/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
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
				fieldname.Token: config.Token,
			},
		})
	})
}

type Config struct {
	Token string `json:"token"`
}

func todoistConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Token: in.ItemFields[fieldname.Token],
	}
	contents, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}
