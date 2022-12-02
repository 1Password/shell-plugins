package gitlab

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
		DocsURL:       sdk.URL("https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html"),
		ManagementURL: sdk.URL("https://gitlab.com/-/profile/personal_access_tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to GitLab.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 26,
					Prefix: "glpat-",
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "URL where GitLab is hosted. Defaults to 'https://gitlab.com'.",
				Optional:            true,
			},
			{
				Name:                fieldname.APIHost,
				MarkdownDescription: "URL where the GitLab API is hosted.",
				Optional:            true,
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryEnvVarPair(defaultEnvVarMapping),
			TryGlabConfigFile(),
		),
	}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"GITLAB_TOKEN":    fieldname.Token,
	"GITLAB_HOST":     fieldname.Host,
	"GITLAB_API_HOST": fieldname.APIHost,
}

func TryGlabConfigFile() sdk.Importer {
	return importer.TryFile("~/.config/glab-cli/config.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config GlabConfig
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		for hostname, hostConfig := range config.Hosts {
			if hostConfig.Token == "" {
				continue
			}

			fields := map[sdk.FieldName]string{
				fieldname.Token: hostConfig.Token,
			}

			nameHint := ""
			if hostname != "gitlab.com" {
				fields[fieldname.Host] = hostname
				nameHint = hostname
			}

			if hostConfig.APIHost != "" && hostConfig.APIHost != "gitlab.com" {
				fields[fieldname.APIHost] = hostConfig.APIHost
			}

			out.AddCandidate(sdk.ImportCandidate{
				Fields:   fields,
				NameHint: nameHint,
			})
		}
	})
}

type GlabConfig struct {
	Hosts map[string]GlabHost `yaml:"hosts"`
}

type GlabHost struct {
	Token   string `yaml:"token"`
	APIHost string `yaml:"api_host"`
}
