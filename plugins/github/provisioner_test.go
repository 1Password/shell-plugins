package github

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestGitHubCLIProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, GitHubCLIProvisioner{}, map[string]plugintest.ProvisionCase{
		"github.com": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GH_TOKEN":     "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				},
			},
		},
		"github.com with explicit host": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				fieldname.Host:  "github.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GH_TOKEN":     "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				},
			},
		},
		"GitHub Enterprise Server": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				fieldname.Host:  "enterprise.github.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GH_ENTERPRISE_TOKEN":     "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GITHUB_ENTERPRISE_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GH_HOST":                 "enterprise.github.com",
				},
			},
		},
		"GitHub Enterprise Cloud": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				fieldname.Host:  "company.ghe.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GH_TOKEN":     "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GH_HOST":      "company.ghe.com",
				},
			},
		},
		"github.localhost": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
				fieldname.Host:  "github.localhost",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GH_TOKEN":     "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GITHUB_TOKEN": "github_pat_OYXGsaLFxgNy9msXs44LFNzg3wh0VsXRGycViVc0iKPOqczc1QKlB3ZVVrm5ESukqKR8nE3jzPBEXAMPLE",
					"GH_HOST":      "github.localhost",
				},
			},
		},
	})
}
