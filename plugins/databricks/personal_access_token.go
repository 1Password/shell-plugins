package databricks

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
		Name:    credname.PersonalAccessToken,
		DocsURL: sdk.URL("https://docs.databricks.com/dev-tools/auth.html"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The host URL of your Databricks instance. Should start with https://",
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Uppercase: true,
						Digits:    true,
						Symbols:   true,
					},
					Prefix: "https://",
				},
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Personal access token used to authenticate to Databricks.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 38,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
						Specific:  []rune{'-'},
					},
					Prefix: "dapi",
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryDatabricksConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"DATABRICKS_HOST":  fieldname.Host,
	"DATABRICKS_TOKEN": fieldname.Token,
}

func TryDatabricksConfigFile() sdk.Importer {
	return importer.TryFile("~/.databrickscfg", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		for _, section := range credentialsFile.Sections() {
			profileName := section.Name()
			fields := make(map[sdk.FieldName]string)

			if section.HasKey("host") && section.Key("host").Value() != "" {
				fields[fieldname.Host] = section.Key("host").Value()
			}
			if section.HasKey("token") && section.Key("token").Value() != "" {
				fields[fieldname.Token] = section.Key("token").Value()
			}

			// add only candidates with required credential fields
			if fields[fieldname.Host] != "" && fields[fieldname.Token] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields:   fields,
					NameHint: importer.SanitizeNameHint(profileName),
				})
			}
		}
	})
}

type Config struct {
	Host  string
	Token string
}
