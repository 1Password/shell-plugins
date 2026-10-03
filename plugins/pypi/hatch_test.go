package pypi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestHatchCLIProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, HatchCLI().Uses[0].Provisioner, map[string]plugintest.ProvisionCase{
		"sets HATCH_INDEX_USER and HATCH_INDEX_AUTH": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "pypi-hatch789xyz",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"HATCH_INDEX_USER": "__token__",
					"HATCH_INDEX_AUTH": "pypi-hatch789xyz",
				},
			},
		},
	})
}

func TestHatchCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, HatchCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
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
