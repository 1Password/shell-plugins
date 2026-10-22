package buildkite

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
		DocsURL:       sdk.URL("https://buildkite.com/docs/platform/cli"),
		ManagementURL: sdk.URL("https://buildkite.com/user/api-access-tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "Organization slug for your Buildkite account.",
				Secret:              false,
				Composition: &schema.ValueComposition{
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
						Symbols:   false,
					},
				},
			},
			{
				Name:                fieldname.Token,
				MarkdownDescription: "API Token used to authenticate with Buildkite.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 45,
					Prefix: "bkua_",
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
			TryBuildkiteConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"BUILDKITE_ORGANIZATION_SLUG": fieldname.Organization,
	"BUILDKITE_API_TOKEN": fieldname.Token,
}

// Check if the platform stores the API Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryBuildkiteConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/bk.yaml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if len(config.Organizations) == 0 {
			return
		}

		for organizationSlug, organization := range config.Organizations {

			if organizationSlug == "" || organization.Token == "" {
				return
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: organization.Token,
					fieldname.Organization: organizationSlug,
				},
				NameHint: importer.SanitizeNameHint(organizationSlug),
			})
		}
	})
}

type Config struct {
    Organizations map[string]Organization `yaml:"organizations"`
}

type Organization struct {
    Token string `yaml:"api_token"`
}
