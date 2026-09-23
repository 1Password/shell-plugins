package pypi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestTwineCLIProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, TwineCLI().Uses[0].Provisioner, map[string]plugintest.ProvisionCase{
		"sets TWINE_USERNAME and TWINE_PASSWORD": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "pypi-AgEIcHlwaS5vcmc",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"TWINE_USERNAME": "__token__",
					"TWINE_PASSWORD": "pypi-AgEIcHlwaS5vcmc",
				},
			},
		},
	})
}

func TestTwineCLINeedsAuth(t *testing.T) {
	plugintest.TestNeedsAuth(t, TwineCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"requires auth for upload": {
			Args:              []string{"upload", "dist/*"},
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
		"skips auth for check command": {
			Args:              []string{"check", "dist/*"},
			ExpectedNeedsAuth: false,
		},
	})
}
