package sourcegraph

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Endpoint: "https://sourcegraph.com",
				fieldname.Token:    "sgp_fake0123456789ab_fake0123456789abcdef0123456789abcdef0123",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SRC_ENDPOINT":     "https://sourcegraph.com",
					"SRC_ACCESS_TOKEN": "sgp_fake0123456789ab_fake0123456789abcdef0123456789abcdef0123",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"SRC_ENDPOINT":     "https://sourcegraph.com",
				"SRC_ACCESS_TOKEN": "sgp_fake0123456789ab_fake0123456789abcdef0123456789abcdef0123",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Endpoint: "https://sourcegraph.com",
						fieldname.Token:    "sgp_fake0123456789ab_fake0123456789abcdef0123456789abcdef0123",
					},
				},
			},
		},
	})
}
