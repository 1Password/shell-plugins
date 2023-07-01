package contentful

import (
	"context"
	"encoding/json"

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
		DocsURL: sdk.URL("https://www.contentful.com/developers/docs/references/content-management-api/#/reference/personal-access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Personal Access Token used to authenticate to Contentful.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 49,
					Prefix: "CFPAT-",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.TempFile(
			contentfulConfig,
			provision.AtFixedPath("~/.contentfulrc.json"),
		),
		Importer: importer.TryAll(
			TryContentfulConfigFile(),
		)}
}

func TryContentfulConfigFile() sdk.Importer {
	return importer.TryFile("~/.contentfulrc.json", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToJSON(&config); err != nil {
			out.AddError(err)
			return
		}

		if config.ManagementToken == "" {
			return
		}

		out.AddCandidate(sdk.ImportCandidate{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: config.ManagementToken,
			},
		})
	})
}

type Config struct {
	ManagementToken     string `json:"managementToken"`
	ActiveEnvironmentId string `json:"activeEnvironmentId"`
	Host                string `json:"host"`
}

func contentfulConfig(in sdk.ProvisionInput) ([]byte, error) {
	config := Config{
		ManagementToken:     in.ItemFields[fieldname.Token],
		ActiveEnvironmentId: "master",
		Host:                "api.contentful.com",
	}
	contents, err := json.MarshalIndent(&config, "", "	")
	if err != nil {
		return nil, err
	}
	return []byte(contents), nil
}
