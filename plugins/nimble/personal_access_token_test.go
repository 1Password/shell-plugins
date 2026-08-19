package nimble

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "ghp_ahhfw82h48fh72nfn29fn291nwidhf8EXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NIMBLE_GITHUB_API_TOKEN": "ghp_ahhfw82h48fh72nfn29fn291nwidhf8EXAMPLE",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"NIMBLE_GITHUB_API_TOKEN": "ghp_ahhfw82h48fh72nfn29fn291nwidhf8EXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "ghp_ahhfw82h48fh72nfn29fn291nwidhf8EXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.nimble/github_api_token": plugintest.LoadFixture(t, "github_api_token"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "ghp_ahhfw82h48fh72nfn29fn291nwidhf8EXAMPLE",
					},
				},
			},
		},
	})
}
