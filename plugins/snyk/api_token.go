package snyk

import (
	"context"
	"os"
	"strings"

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
		DocsURL:       sdk.URL("https://docs.snyk.io/snyk-api-info/authentication-for-api"),
		ManagementURL: sdk.URL("https://app.snyk.io/account"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Snyk.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 36,
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
			TrySnykConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"SNYK_TOKEN": fieldname.Token,
}

type Config struct {
	Token string `json:"api"`
}

func TrySnykConfigFile() sdk.Importer {
	path := configFilePath()
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
				fieldname.Token: config.Token,
			},
		})
	})
}

func configFilePath() string {
	configPath := os.Getenv("XDG_CONFIG_HOME")
	if configPath == "" {
		configPath = "~/.config"
	}

	if !strings.HasSuffix(configPath, "/") {
		configPath += "/"
	}

	configPath += "configstore/snyk.json"
	return configPath
}
