package apono

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(
		t, PersonalAPIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
			"default": {
				ItemFields: map[sdk.FieldName]string{
					fieldname.Token: "apono_example_personal_token_ynAvsvBcxSGzMkPr",
				},
				ExpectedOutput: sdk.ProvisionOutput{
					CommandLine: []string{"--personal-token", "apono_example_personal_token_ynAvsvBcxSGzMkPr"},
				},
			},
		},
	)
}

func TestPersonalAPITokenImporter(t *testing.T) {
	expectedCandidates := []sdk.ImportCandidate{
		{
			Fields: map[sdk.FieldName]string{
				fieldname.Token: "apono_example_personal_token_ynAvsvBcxSGzMkPr",
			},
			NameHint: "work",
		},
	}

	plugintest.TestImporter(
		t, PersonalAPIToken().Importer, map[string]plugintest.ImportCase{
			"config file (macOS)": {
				OS: "darwin",
				Files: map[string]string{
					"~/Library/Application Support/apono-cli/config.json": plugintest.LoadFixture(t, "config.json"),
				},
				ExpectedCandidates: expectedCandidates,
			},
			"config file (Linux)": {
				OS: "linux",
				Files: map[string]string{
					"~/.config/apono-cli/config.json": plugintest.LoadFixture(t, "config.json"),
				},
				ExpectedCandidates: expectedCandidates,
			},
		},
	)
}
