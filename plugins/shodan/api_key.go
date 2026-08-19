package shodan

import (
	"context"

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
		DocsURL:       sdk.URL("https://cli.shodan.io"),
		ManagementURL: sdk.URL("https://account.shodan.io"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to Shodan.",
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
		DefaultProvisioner: provision.TempFile(
			provision.FieldAsFile(fieldname.APIKey),
			provision.AtFixedPath("~/.shodan/api_key"),
		),
		Importer: importer.TryAll(
			TryShodanConfigFile("~/.shodan/api_key"),
			TryShodanConfigFile("~/.config/shodan/api_key"),
		)}
}

func TryShodanConfigFile(configPath string) sdk.Importer {
	return importer.TryFile(configPath, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		apiKey := contents.ToString()

		if apiKey == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: apiKey,
			},
		})
	})
}
