package fastly

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
		DocsURL:       sdk.URL("https://docs.fastly.com/en/guides/using-api-tokens"),
		ManagementURL: sdk.URL("https://manage.fastly.com/account/personal/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API Token used to authenticate to Fastly.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 32,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "FASTLY_API_TOKEN"),
			importer.MacOnly(
				TryFastlyConfigFile("~/Library/Application Support/fastly/config.toml"),
			),
			importer.LinuxOnly(
				TryFastlyConfigFile("~/.config/fastly/config.toml"),
			),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"FASTLY_API_TOKEN": fieldname.Token,
}

func TryFastlyConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		for profileName, configProfile := range config.Profile {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: configProfile.Token,
				},
				NameHint: importer.SanitizeNameHint(profileName),
			})
		}
	})
}

type ConfigProfile struct {
	Email string `toml:"email"`
	Token string `toml:"token"`
}

type Config struct {
	Profile map[string]ConfigProfile
}
