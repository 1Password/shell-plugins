package nirmata

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
		DocsURL:       sdk.URL("https://nirmata.com/docs/api_token"),               // TODO: Replace with actual URL
		ManagementURL: sdk.URL("https://console.nirmata.com/user/security/tokens"), // TODO: Replace with actual URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Nirmata.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 116,
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Email,
				MarkdownDescription: "Email address registered in Nirmata.",
			},
			{
				Name:                fieldname.Address,
				MarkdownDescription: "Url address of Nirmata[https://nirmata.io].",
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),

		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryNirmataConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"NIRMATA_TOKEN": fieldname.Token,
	"NIRMATA_URL":   fieldname.URL,
}

func TryNirmataConfigFile() sdk.Importer {
	return importer.TryFile("~/.nirmata/config", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {

		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}
		for _, section := range credentialsFile.Sections() {
			fields := make(map[sdk.FieldName]string)
			if section.HasKey("address") && section.Key("address").Value() != "" {
				fields[fieldname.Address] = section.Key("address").Value()
			}
			if section.HasKey("email") && section.Key("email").Value() != "" {
				fields[fieldname.Email] = section.Key("email").Value()
			}
			if section.HasKey("token") && section.Key("token").Value() != "" {
				fields[fieldname.Token] = section.Key("token").Value()
			}
			if fields[fieldname.Address] != "" && fields[fieldname.Email] != "" && fields[fieldname.Token] != "" {
				out.AddCandidate(sdk.ImportCandidate{
					Fields: fields,
				})
			}

		}

	})
}

type Config struct {
	Address string
	Email   string
	Token   string
}
