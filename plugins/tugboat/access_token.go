package tugboat

import (
	"context"
	"fmt"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessToken,
		DocsURL:       sdk.URL("https://docs.tugboatqa.com/tugboat-cli/set-an-access-token/"),
		ManagementURL: sdk.URL("https://dashboard.tugboatqa.com/access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Tugboat.",
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
		DefaultProvisioner: provision.TempFile(tugboatConfig, provision.AtFixedPath("~/.tugboat.yml")),
		Importer: importer.TryAll(
			TryTugboatConfigFile(),
		)}
}

func tugboatConfig(in sdk.ProvisionInput) ([]byte, error) {
	content := ""

	if token, ok := in.ItemFields[fieldname.Token]; ok {
		content += configFileEntry("token", token)
	}

	return []byte(content), nil
}

func configFileEntry(key string, value string) string {
	return fmt.Sprintf("%s: %s\n", key, value)
}

func TryTugboatConfigFile() sdk.Importer {
	return importer.TryFile("~/.tugboat.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
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

type Config struct {
	Token string
}
