package tunnelto

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
		ManagementURL: sdk.URL("https://dashboard.tunnelto.dev/"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "API Key used to authenticate to tunnelto.dev.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 22,
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
			provision.AtFixedPath("~/.tunnelto/key.token"),
		),
		Importer: importer.TryAll(
			TrytunneltodevConfigFile(),
		)}
}

func TrytunneltodevConfigFile() sdk.Importer {
	return importer.TryFile("~/.tunnelto/key.token", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		if contents.ToString() == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.APIKey: contents.ToString(),
			},
		})
	})
}
