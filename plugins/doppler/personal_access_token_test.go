package doppler

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"DOPPLER_TOKEN environment variable": {
			Environment: map[string]string{
				"DOPPLER_TOKEN": "dp.pt.SQgRDoLc2lBYVu5Vr2T4XHPvBcp0HlhMZq8F11whbvQEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dp.pt.SQgRDoLc2lBYVu5Vr2T4XHPvBcp0HlhMZq8F11whbvQEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.doppler/.doppler.yaml": plugintest.LoadFixture(t, "doppler.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dp.ct.0nFkB2tJ8wQxZ5pYcR7vLmA3sD6gH9jKqW1eU4iO0EXAMPLE",
					},
					NameHint: "example-project",
				},
			},
		},
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "dp.pt.SQgRDoLc2lBYVu5Vr2T4XHPvBcp0HlhMZq8F11whbvQEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DOPPLER_TOKEN": "dp.pt.SQgRDoLc2lBYVu5Vr2T4XHPvBcp0HlhMZq8F11whbvQEXAMPLE",
				},
			},
		},
	})
}
