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
					Fields: map[string]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
				{
					Fields: map[string]string{
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
					Fields: map[string]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
				{
					Fields: map[string]string{
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
					Fields: map[string]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GH_ENTERPRISE_TOKEN": {
			Environment: map[string]string{
				"GH_ENTERPRISE_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GITHUB_ENTERPRISE_TOKEN": {
			Environment: map[string]string{
				"GITHUB_ENTERPRISE_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					},
				},
			},
		},
		"GH_HOST": {
			Environment: map[string]string{
				"GH_HOST": "host",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Host: "host",
					},
				},
				{
					Fields: map[string]string{
						fieldname.Host: "host",
					},
				},
				{
					Fields: map[string]string{
						fieldname.Host: "host",
					},
				},
				{
					Fields: map[string]string{
						fieldname.Host: "host",
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().Provisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[string]string{
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
