package tunnelto

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"temp file": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.APIKey: "XddpK7jZiQ0CpE3EXAMPLE",
			},
			CommandLine: []string{"tunnelto"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"tunnelto"},
				Files: map[string]sdk.OutputFile{
					"~/.tunnelto/key.token": {
						Contents: []byte(plugintest.LoadFixture(t, "key.token")),
					},
				},
			},
		},
	})
}
