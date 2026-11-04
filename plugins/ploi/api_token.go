package ploi

import (
	"context"
	"regexp"

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
		DocsURL:       sdk.URL("https://developers.ploi.io/?documentation=cli"),
		ManagementURL: sdk.URL("https://ploi.io/profile/api-keys"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Ploi CLI.",
				Secret:              true,
				Composition: &schema.ValueComposition{
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
			TryPloiCLIConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"PLOI_API_TOKEN": fieldname.Token,
}

func TryPloiCLIConfigFile() sdk.Importer {
	return importer.TryFile("~/.ploi/config.php", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		tokenRegexp := regexp.MustCompile(`'token'\s*=>\s*'([^']*)'`)
		matches := tokenRegexp.FindStringSubmatch(contents.ToString())
		if len(matches) < 2 || matches[1] == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: matches[1],
			},
		})
	})
}
