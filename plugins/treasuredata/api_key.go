package treasuredata

import (
	"context"
	"os"

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
		DocsURL:       sdk.URL("https://docs.treasuredata.com/display/public/PD/Configuring+Authentication+for+TD+Using+the+TD+Toolbelt"),
		ManagementURL: sdk.URL("https://console.treasuredata.com/app/mp/ak"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Treasure Data.",
				Secret:              true,
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
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"TREASURE_DATA_API_KEY": fieldname.APIKey,
			}),
			TryTreasureDataConfigFile(os.Getenv("TREASURE_DATA_CONFIG_PATH")),
			TryTreasureDataConfigFile(os.Getenv("TD_CONFIG_PATH")),
			TryTreasureDataConfigFile("~/.td/td.conf"),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"TD_API_KEY": fieldname.APIKey,
}

func TryTreasureDataConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("apikey") && section.Key("apikey").Value() != "" {
				fields[fieldname.APIKey] = section.Key("apikey").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}
