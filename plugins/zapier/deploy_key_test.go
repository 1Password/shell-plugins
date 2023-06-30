package zapier

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestDeployKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, DeployKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"config file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Key: "bd0moxtb5vk3koppai1f3ituiexample",
			},
			CommandLine: []string{"zapier"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"zapier"},
				Files: map[string]sdk.OutputFile{
					"~/.zapierrc": {
						Contents: []byte(plugintest.LoadFixture(t, ".zapierrc")),
					},
				},
			},
		},
	})
}

func TestDeployKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, DeployKey().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.zapierrc": plugintest.LoadFixture(t, ".zapierrc"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Key: "bd0moxtb5vk3koppai1f3ituiexample",
					},
				},
			},
		},
	})
}
