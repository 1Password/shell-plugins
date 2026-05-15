package zapier

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDeployKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DeployKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"environment": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key: "g22y9yajfqqcogilnrnm5lxv6example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"ZAPIER_DEPLOY_KEY": "g22y9yajfqqcogilnrnm5lxv6example",
				},
			},
		},
	})
}

func TestDeployKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, DeployKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"ZAPIER_DEPLOY_KEY": "g22y9yajfqqcogilnrnm5lxv6example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "g22y9yajfqqcogilnrnm5lxv6example",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.zapierrc": plugintest.LoadFixture(t, ".zapierrc"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "g22y9yajfqqcogilnrnm5lxv6example",
					},
				},
			},
		},
	})
}
