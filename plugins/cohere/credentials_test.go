package cohere

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCredentialsProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.JWT:   "PLEzI1N--EXAMPLE",
				fieldname.URL:   "https://api.os.cohere.ai",
				fieldname.Email: "wendy@appleseed.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					"~/.command/config": {
						Contents: []byte(plugintest.LoadFixture(t, "config")),
					},
				},
			},
		},
	})
}

func TestCredentialsImporter(t *testing.T) {
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{

		"config file": {
			Files: map[string]string{
				"~/.command/config": plugintest.LoadFixture(t, "config"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.JWT:   "PLEzI1N--EXAMPLE",
						fieldname.URL:   "https://api.os.cohere.ai",
						fieldname.Email: "wendy@appleseed.com",
					},
				},
			},
		},
	})
}
