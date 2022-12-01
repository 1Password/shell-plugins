package datadog

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"DATADOG_API_KEY": "vvt4tzq3td4cb3wklna5bclv0example",
				"DATADOG_APP_KEY": "g3jesabd5z8mm61ax7a0k8rzgju1yeps4example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.APIKey: "vvt4tzq3td4cb3wklna5bclv0example",
						fieldname.AppKey: "g3jesabd5z8mm61ax7a0k8rzgju1yeps4example",
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
				fieldname.APIKey: "vvt4tzq3td4cb3wklna5bclv0example",
				fieldname.AppKey: "g3jesabd5z8mm61ax7a0k8rzgju1yeps4example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"DATADOG_API_KEY": "vvt4tzq3td4cb3wklna5bclv0example",
					"DATADOG_APP_KEY": "g3jesabd5z8mm61ax7a0k8rzgju1yeps4example",
				},
			},
		},
	})
}
