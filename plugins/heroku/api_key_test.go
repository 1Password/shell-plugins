package heroku

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPIKeyImporter(t *testing.T) {
	plugintest.TestImporter(t, APIKey().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"HEROKU_API_KEY": "dh7k7m662pqglxaybr1p0gpg1cu33example",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.APIKey: "dh7k7m662pqglxaybr1p0gpg1cu33example",
					},
				},
			},
		},
		"netrc file": {
			Files: map[string]string{
				"~/.netrc": plugintest.LoadFixture(t, "netrc"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "wendy@appleseed.com",
					Fields: map[string]string{
						fieldname.APIKey: "dh7k7m662pqglxaybr1p0gpg1cu33example",
					},
				},
				{
					NameHint: "wendy@appleseed.com",
					Fields: map[string]string{
						fieldname.APIKey: "dh7k7m662pqglxaybr1p0gpg1cu33example",
					},
				},
			},
		},
		"netrc file non-Heroku": {
			Files: map[string]string{
				"~/.netrc": plugintest.LoadFixture(t, "netrc-non-heroku"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
	})
}

func TestAPIKeyProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIKey().Provisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[string]string{
				fieldname.APIKey: "dh7k7m662pqglxaybr1p0gpg1cu33example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"HEROKU_API_KEY": "dh7k7m662pqglxaybr1p0gpg1cu33example",
				},
			},
		},
	})
}
