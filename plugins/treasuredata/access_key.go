package treasuredata

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AccessKey() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AccessKey,
		DocsURL:       sdk.URL("https://docs.treasuredata.com/display/public/PD/TD+Toolbelt"),
		ManagementURL: sdk.URL("https://console.treasuredata.com/app/mp/ak"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.User,
				MarkdownDescription: "User name specified by email",
			},
			{
				Name:                fieldname.APIKey,
				MarkdownDescription: "APIKey used to authenticate to Treasure Data.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			TryTreasureDataConfigFile("~/.td/td.conf"),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"TD_API_KEY": fieldname.APIKey,
}

func TryTreasureDataConfigFile(path string) sdk.Importer {
	return importer.TryFile(path, func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		credentialsFile, err := contents.ToINI()
		if err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)
		for _, section := range credentialsFile.Sections() {
			if section.HasKey("user") && section.Key("user").Value() != "" {
				fields[fieldname.User] = section.Key("user").Value()
			}
			if section.HasKey("apikey") && section.Key("apikey").Value() != "" {
				fields[fieldname.APIKey] = section.Key("apikey").Value()
			}
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: fields,
		})
	})
}
