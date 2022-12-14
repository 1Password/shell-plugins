package linode

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"LINODE_CLI_TOKEN": "cn5z4umbimz0lxgzvps1bl979n8lpwnu6qmb4x19bddzx6siormnoxg2yexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "cn5z4umbimz0lxgzvps1bl979n8lpwnu6qmb4x19bddzx6siormnoxg2yexample",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.config/linode-cli": plugintest.LoadFixture(t, "linode-cli"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "cn5z4umbimz0lxgzvps1bl979n8lpwnu6qmb4x19bddzx6siormnoxg2yexample",
					},
				},
			},
		},
	})
}

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "cn5z4umbimz0lxgzvps1bl979n8lpwnu6qmb4x19bddzx6siormnoxg2yexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"LINODE_CLI_TOKEN": "cn5z4umbimz0lxgzvps1bl979n8lpwnu6qmb4x19bddzx6siormnoxg2yexample",
				},
			},
		},
	})
}
