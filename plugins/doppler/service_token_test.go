package doppler

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestServiceTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, ServiceToken().Importer, map[string]plugintest.ImportCase{
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
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "dp.pt.SQgRDoLc2lBYVu5Vr2T4XHPvBcp0HlhMZq8F11whbvQEXAMPLE",
					},
				},
			},
		},
	})
}

func TestServiceTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, ServiceToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
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
