package github

import (
	"context"
	"strings"

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
		DocsURL:       sdk.URL("https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token"),
		ManagementURL: sdk.URL("https://github.com/settings/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to GitHub.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "ghp_",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The GitHub host to authenticate to. Defaults to 'github.com'.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			importer.TryAllEnvVars(fieldname.Token, "GH_TOKEN", "GITHUB_PAT"),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":             fieldname.Host,
				"GH_ENTERPRISE_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":                 fieldname.Host,
				"GITHUB_ENTERPRISE_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":  fieldname.Host,
				"GH_TOKEN": fieldname.Token,
			}),
			importer.TryEnvVarPair(map[string]sdk.FieldName{
				"GH_HOST":      fieldname.Host,
				"GITHUB_TOKEN": fieldname.Token,
			}),
			TryGitHubConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"GH_TOKEN": fieldname.Token,
	"GH_ENTERPRISE_TOKEN": fieldname.Token,
        "GH_HOST": fieldname.Host,
}

type Config struct {
	Token string `yaml:"oauth_token"`
}

func TryGitHubConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/gh/hosts.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config map[string]Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for host, values := range config {
			if strings.HasPrefix(values.Token, "ghp_") {
				candidate := sdk.ImportCandidate{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: values.Token,
					},
				}

				if host != "github.com" {
					candidate.NameHint = host
					candidate.Fields[fieldname.Host] = host
				}

				out.AddCandidate(candidate)
			}
		}
	})
}
