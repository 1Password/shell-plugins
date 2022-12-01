package snyk

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"default": {
			Environment: map[string]string{
				"SNYK_TOKEN": "etacgrrwj86t255ckguircl3kw3ftexample",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[string]string{
						fieldname.Token: "etacgrrwj86t255ckguircl3kw3ftexample",
					},
				},
			},
		},
	})
}

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().Provisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[string]string{
				fieldname.Token: "etacgrrwj86t255ckguircl3kw3ftexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SNYK_TOKEN": "etacgrrwj86t255ckguircl3kw3ftexample",
				},
			},
		},
	})
}
