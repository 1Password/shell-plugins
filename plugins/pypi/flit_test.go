package pypi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestFlitCLIProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, FlitCLI().Uses[0].Provisioner, map[string]plugintest.ProvisionCase{
		"sets FLIT_USERNAME and FLIT_PASSWORD": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "pypi-flit123abc",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"FLIT_USERNAME": "__token__",
					"FLIT_PASSWORD": "pypi-flit123abc",
				},
			},
		},
	})
}

func TestFlitCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, FlitCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"requires auth for publish": {
			Args:              []string{"publish"},
			ExpectedNeedsAuth: true,
		},
		"skips auth for help": {
			Args:              []string{"--help"},
			ExpectedNeedsAuth: false,
		},
		"skips auth for version": {
			Args:              []string{"--version"},
			ExpectedNeedsAuth: false,
		},
		"skips auth for build command": {
			Args:              []string{"build"},
			ExpectedNeedsAuth: false,
		},
	})
}
