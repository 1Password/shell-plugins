package flyctl

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
				fieldname.Token: "DtP7HoOPOBHMP6bE5tx3nguB5r2zPpSbg9hlEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"FLY_ACCESS_TOKEN": "DtP7HoOPOBHMP6bE5tx3nguB5r2zPpSbg9hlEXAMPLE",
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"FLY_ACCESS_TOKEN": "DtP7HoOPOBHMP6bE5tx3nguB5r2zPpSbg9hlEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "DtP7HoOPOBHMP6bE5tx3nguB5r2zPpSbg9hlEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.fly/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "DtP7HoOPOBHMP6bE5tx3nguB5r2zPpSbg9hlEXAMPLE",
					},
				},
			},
		},
	})
}
