package gitlab

import (
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

var defaultEnvVarMapping = map[string]string{
	fieldname.Token:   "GITLAB_TOKEN",
	fieldname.Host:    "GITLAB_HOST",
	fieldname.APIHost: "GITLAB_API_HOST",
}

func TryGlabConfigFile() sdk.Importer {
	// TODO: Try importing token from ~/.config/glab-cli/config.yml
	return importer.NoOp()
}
