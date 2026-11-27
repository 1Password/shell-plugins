package vercel

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
		DocsURL:       sdk.URL("https://vercel.com/docs/rest-api#authentication"),
		ManagementURL: sdk.URL("https://vercel.com/account/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Vercel.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 24,
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.MacOnly(TryVercelConfigFile("~/Library/Application Support/com.vercel.cli/auth.json")),
			importer.LinuxOnly(TryVercelConfigFile("~/.config/com.vercel.cli/auth.json")),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"VERCEL_TOKEN": fieldname.Token,
}

func TryVercelConfigFile(path string) sdk.Importer {
	return importer.TryFile(
		path,
		func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
			var config Config
			if err := contents.ToJSON(&config); err != nil {
				out.AddError(err)
				return
			}

			if config.Token == "" {
				return
			}

			out.AddCandidate(
				sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: config.Token,
					},
				},
			)
		},
	)
}

type Config struct {
	Token string `json:"token"`
}
