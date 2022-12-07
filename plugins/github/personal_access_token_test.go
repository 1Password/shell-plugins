package github

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"GITHUB_TOKEN": {
			Environment: map[string]string{
				"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GH_TOKEN": {
			Environment: map[string]string{
				"GH_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GITHUB_PAT": {
			Environment: map[string]string{
				"GITHUB_PAT": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GitHub Enterprise": {
			Environment: map[string]string{
				"GH_ENTERPRISE_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				"GH_HOST":             "github.acme.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
						fieldname.Host:  "github.acme.com",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host: "github.acme.com",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host: "github.acme.com",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Host: "github.acme.com",
					},
				},
			},
		},
		"GitHub config file with github_pat token prefix": {
			Files: map[string]string{
				"~/.config/gh/hosts.yml": plugintest.LoadFixture(t, "hosts.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_8BQUJmcVkoNo48YBxXjyF20gqFi7hYLdzhGmJUISVUCjyHOA0sdeO8Xmw7LqYY0Ng1wndPbww6fEXAMPLE",
					},
				},
				{
					NameHint: "enterprise.github.com",
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "github_pat_2BaueV0i4Jv832lprA6lLavL8H2cw9IPBVKqdmHdqKrYstzQUwTukwpX54Z8HAN3ehKGid6poaiEXAMPLE",
						fieldname.Host:  "enterprise.github.com",
					},
				},
			},
		},
		"GitHub config file with gho token prefix": {
			Files: map[string]string{
				"~/.config/gh/hosts.yml": plugintest.LoadFixture(t, "hosts_gho_type.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				},
			},
		},
	})
}
