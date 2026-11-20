package prefect

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.prefect.io/cloud/users/api-keys/"),
		ManagementURL: sdk.URL("https://app.prefect.cloud/my/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "Key used to authenticate to Prefect.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "pnu_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.URL,
				MarkdownDescription: "URL for the Prefect workspace.",
				Secret:              false,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
		)}
}

func TryPrefectConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		for profileName, configProfile := range config.Profile {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.APIKey: configProfile.APIKey,
					fieldname.URL:    configProfile.URL,
				},
				NameHint: importer.SanitizeNameHint(profileName),
			})
		}
	})
}

type ConfigProfile struct {
	URL    string `toml:"PREFECT_API_URL"`
	APIKey string `toml:"PREFECT_API_KEY"`
}

type Config struct {
	Profile map[string]ConfigProfile
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"PREFECT_API_KEY": fieldname.APIKey,
	"PREFECT_API_URL": fieldname.URL,
}
