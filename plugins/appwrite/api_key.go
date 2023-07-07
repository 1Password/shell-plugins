package appwrite

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

func ConfigPath() string {
	return "~/.appwrite/prefs.json"
}

func APIKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.APIKey,
		DocsURL:       sdk.URL("https://appwrite.io/docs"),
		ManagementURL: sdk.URL("https://cloud.appwrite.io/console/account"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Appwrite.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 256,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Endpoint,
				MarkdownDescription: "Appwrite server endpoint.",
				Secret:              false,
				Optional:            false,
			},
		},
		DefaultProvisioner: provision.TempFile(appwriteConfig, provision.AtFixedPath(ConfigPath())),
		Importer:           TryAppwriteConfigFile(),
	}
}

func appwriteConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		APIKey:   in.ItemFields[fieldname.APIKey],
		Endpoint: in.ItemFields[fieldname.Endpoint],
	}

	contents, err := json.Marshal(&config)
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}

func TryAppwriteConfigFile() sdk.Importer {
	return importer.TryFile(ConfigPath(), func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.APIKey == "" {
			return
		}

		if config.Endpoint == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey:   config.APIKey,
				fieldname.Endpoint: config.Endpoint,
			},
		})
	})
}

type Config struct {
	APIKey   string `json:"key"`
	Endpoint string `json:"endpoint"`
}
