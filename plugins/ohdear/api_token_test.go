package ohdear

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "SZ5rluwzbtMyyQFQNoeqEFbpVbTL0ItsXEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"OHDEAR_CONFIG_FILE_PATH": "/tmp/config.json",
				},
				Files: map[string]sdk.OutputFile{
					"/tmp/config.json": {
						Contents: []byte(plugintest.LoadFixture(t, "import_or_provision")),
					},
				},
			},
		},
	})
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"config file": {
			Files: map[string]string{
				"~/.ohdear/config.json": plugintest.LoadFixture(t, "import_or_provision"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "SZ5rluwzbtMyyQFQNoeqEFbpVbTL0ItsXEXAMPLE",
					},
				},
			},
		},
	})
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"OHDEAR_API_TOKEN": "SZ5rluwzbtMyyQFQNoeqEFbpVbTL0ItsXEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "SZ5rluwzbtMyyQFQNoeqEFbpVbTL0ItsXEXAMPLE",
					},
				},
			},
		},
	})
}
