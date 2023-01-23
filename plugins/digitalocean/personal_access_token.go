package digitalocean

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func PersonalAccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.PersonalAccessToken,
		DocsURL:       sdk.URL("https://docs.digitalocean.com/reference/api/create-personal-access-token/"),
		ManagementURL: sdk.URL("https://cloud.digitalocean.com/account/api/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to DigitalOcean.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 71,
					Prefix: "dop_v1_",
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(map[string]sdk.FieldName{
			"DIGITALOCEAN_ACCESS_TOKEN": fieldname.Token,
		}),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "DIGITALOCEAN_ACCESS_TOKEN"),
			importer.MacOnly(
				TryDigitalOceanConfigFile("~/Library/Application Support/doctl/config.yaml")
			),
		),
	}
}

func TryDigitalOceanConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.AccessToken == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: config.AccessToken,
			},
		})
	})
}

type Config struct {
	AccessToken string `yaml:"access-token"`
}
