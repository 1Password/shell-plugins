package ohdear

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
		DocsURL:       sdk.URL("https://ohdear.app/docs/api/introduction#get-your-api-token"),
		ManagementURL: sdk.URL("https://ohdear.app/user/api-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Oh Dear.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(configFile,
			provision.Filename("config.json"),
			provision.SetPathAsEnvVar("OHDEAR_CONFIG_FILE_PATH"),
		),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOhdearConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OHDEAR_API_TOKEN": fieldname.Token,
}

func configFile(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		Token: in.ItemFields[fieldname.Token],
	}

	contents, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}

func TryOhdearConfigFile() sdk.Importer {
	return importer.TryFile("~/.ohdear/config.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
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
