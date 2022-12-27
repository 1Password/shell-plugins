package hcloud

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
		DocsURL:       sdk.URL("https://github.com/hetznercloud/cli"),
		ManagementURL: sdk.URL("https://console.hetzner.cloud/projects"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Hetzner Cloud.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 64,
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
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryHetznerCloudConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"HCLOUD_TOKEN": fieldname.Token,
}

func TryHetznerCloudConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/hcloud/cli.toml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToTOML(&config); err != nil {
			out.AddError(err)
			return
		}

		for _, configContext := range config.Contexts {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: configContext.Token,
				},
				NameHint: importer.SanitizeNameHint(configContext.Name),
			})
		}
	})
}

type Config struct {
	Contexts []ConfigContext `toml:"contexts"`
}

type ConfigContext struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}
