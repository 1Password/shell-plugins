package gitlab

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"GITLAB_TOKEN": "glpat-sJy3L26ZNW7A3EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token: "glpat-sJy3L26ZNW7A3EXAMPLE",
					},
				},
			},
		},
		"glab config file": {
			Files: map[string]string{
				"~/.config/glab-cli/config.yml": plugintest.LoadFixture(t, "glab-config.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token: "glpat-sJy3L26ZNW7A3EXAMPLE",
					},
				},
			},
		},
		"glab config file with self-hosted instance": {
			Files: map[string]string{
				"~/.config/glab-cli/config.yml": plugintest.LoadFixture(t, "glab-config-self-hosted.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "gitlab.acme.com",
					Fields: map[string]string{
						fieldname.Token:   "glpat-sJy3L26ZNW7A3EXAMPLE",
						fieldname.Host:    "gitlab.acme.com",
						fieldname.APIHost: "api.gitlab.acme.com",
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[string]string{
				fieldname.Token: "glpat-sJy3L26ZNW7A3EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITLAB_TOKEN": "glpat-sJy3L26ZNW7A3EXAMPLE",
				},
			},
		},
		"self-hosted instance": {
			ItemFields: map[string]string{
				fieldname.Token:   "glpat-sJy3L26ZNW7A3EXAMPLE",
				fieldname.Host:    "gitlab.acme.com",
				fieldname.APIHost: "api.gitlab.acme.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITLAB_TOKEN":    "glpat-sJy3L26ZNW7A3EXAMPLE",
					"GITLAB_HOST":     "gitlab.acme.com",
					"GITLAB_API_HOST": "api.gitlab.acme.com",
				},
			},
		},
	})
}
