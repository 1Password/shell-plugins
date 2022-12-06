package okta

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
		DocsURL:       sdk.URL("https://developer.okta.com/docs/guides/create-an-api-token/main/"),
		ManagementURL: nil, // TODO: Add management URL
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to Okta.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 42,
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.OrgURL,
				MarkdownDescription: "URL of the Okta organization to authenticate to.",
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryOktaConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"OKTA_CLIENT_TOKEN":  fieldname.Token,
	"OKTA_CLIENT_ORGURL": fieldname.OrgURL,
}

func TryOktaConfigFile() sdk.Importer {
	return importer.TryFile("~/.okta/okta.yaml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		fields := make(map[sdk.FieldName]string)

		if token := config.Okta.Client.Token; token != "" {
			fields[fieldname.Token] = token
		}

		if orgURL := config.Okta.Client.OrgURL; orgURL != "" {
			fields[fieldname.OrgURL] = orgURL
		}

		if len(fields) > 0 {
			out.AddCandidate(sdk.ImportCandidate{
				Fields: fields,
			})
		}
	})
}

type Config struct {
	Okta OktaConfig `yaml:"okta"`
}

type OktaConfig struct {
	Client ClientConfig `yaml:"client"`
}

type ClientConfig struct {
	OrgURL string `yaml:"orgUrl"`
	Token  string `yaml:"token"`
}
