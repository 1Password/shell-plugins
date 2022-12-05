package circleci

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAPIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"CIRCLECI_CLI_TOKEN": "1evr6rbndnaphymaljwpulrlvws7oolrmexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "1evr6rbndnaphymaljwpulrlvws7oolrmexample",
					},
				},
			},
		},
		"CircleCI config file": {
			Files: map[string]string{
				"~/.circleci/cli.yml": plugintest.LoadFixture(t, "cli.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "1evr6rbndnaphymaljwpulrlvws7oolrmexample",
					},
				},
			},
		},
	})
}

func TestPersonalAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAPIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "1evr6rbndnaphymaljwpulrlvws7oolrmexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"CIRCLECI_CLI_TOKEN": "1evr6rbndnaphymaljwpulrlvws7oolrmexample",
				},
			},
		},
	})
}
