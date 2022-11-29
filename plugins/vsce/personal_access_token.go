package vsce

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
		DocsURL:       sdk.URL("https://code.visualstudio.com/api/working-with-extensions/publishing-extension#vsce"),
		ManagementURL: sdk.URL("https://dev.azure.com"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate with vsce.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 52,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		Provisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryVSCEConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]string{
	fieldname.Token: "VSCE_PAT",
}

type publisher struct {
	pat string
}

type Config struct {
	publishers []publisher
}

func TryVSCEConfigFile() sdk.Importer {
	return importer.TryFile("~/.vsce", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if len(config.publishers) == 0 || len(config.publishers[0].pat) == 0 {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			NameHint: "foo",
			Fields: map[string]string{
				fieldname.Token: config.publishers[0].pat,
			},
		})
	})
}
