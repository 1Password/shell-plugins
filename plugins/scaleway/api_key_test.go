package scaleway

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"Environment variables": {
			Environment: map[string]string{
				"SCW_ACCESS_KEY":              "SCWSYXTFI97NSEXAMPLE",
				"SCW_SECRET_KEY":              "d9b67b48-873c-8ece-8270-e1e15example",
				"SCW_DEFAULT_ORGANIZATION_ID": "11111111-2222-3333-4444-55555example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:         "SCWSYXTFI97NSEXAMPLE",
						fieldname.SecretAccessKey:     "d9b67b48-873c-8ece-8270-e1e15example",
						fieldname.DefaultOrganization: "11111111-2222-3333-4444-55555example",
					},
				},
			},
		},
		"SCW default config file location": {
			Files: map[string]string{
				"~/.config/scw/config.yaml": plugintest.LoadFixture(t, "simple.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:         "SCWSYXTFI97NSEXAMPLE",
						fieldname.SecretAccessKey:     "d9b67b48-873c-8ece-8270-e1e15example",
						fieldname.DefaultOrganization: "11111111-2222-3333-4444-55555example",
					},
				},
			},
		},
		"SCW config file with optional settings": {
			Files: map[string]string{
				"~/.config/scw/config.yaml": plugintest.LoadFixture(t, "optional.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.AccessKeyID:         "SCWSYXTFI97NSEXAMPLE",
						fieldname.SecretAccessKey:     "d9b67b48-873c-8ece-8270-e1e15example",
						fieldname.DefaultOrganization: "11111111-2222-3333-4444-55555example",
						fieldname.DefaultRegion:       "fr-par",
						fieldname.DefaultZone:         "fr-par-1",
					},
				},
			},
		},
	})
}

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.AccessKeyID:         "SCWSYXTFI97NSEXAMPLE",
				fieldname.SecretAccessKey:     "d9b67b48-873c-8ece-8270-e1e15example",
				fieldname.DefaultOrganization: "11111111-2222-3333-4444-55555example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SCW_ACCESS_KEY":              "SCWSYXTFI97NSEXAMPLE",
					"SCW_SECRET_KEY":              "d9b67b48-873c-8ece-8270-e1e15example",
					"SCW_DEFAULT_ORGANIZATION_ID": "11111111-2222-3333-4444-55555example",
				},
			},
		},
	})
}
