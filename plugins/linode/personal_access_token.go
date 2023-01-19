package linode

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
		DocsURL:       sdk.URL("https://www.linode.com/docs/products/tools/cloud-manager/guides/cloud-api-keys/"),
		ManagementURL: sdk.URL("https://cloud.linode.com/profile/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Linode.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 64,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"LINODE_CLI_TOKEN": fieldname.Token,
}

// TryConfigFile looks for the token in the ~/.config/linode-cli file.
func TryConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/linode-cli", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		configFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range configFile.Sections() {
			if section.HasKey("token") && section.Key("token").Value() != "" {
				out.AddCandidate(sdk.ImportCandidate{
					NameHint: importer.SanitizeNameHint(section.Name()),
					Fields: map[sdk.FieldName]string{
						fieldname.Token: section.Key("token").Value(),
					},
				})
			}
		}
	})
}
