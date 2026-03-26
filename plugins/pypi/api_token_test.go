package pypi

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, APIToken().Importer, map[string]plugintest.ImportCase{
		"TWINE_PASSWORD environment variable": {
			Environment: map[string]string{
				"TWINE_PASSWORD": "pypi-AgEIcHlwaS5vcmc",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pypi-AgEIcHlwaS5vcmc",
					},
				},
			},
		},
		"FLIT_PASSWORD environment variable": {
			Environment: map[string]string{
				"FLIT_PASSWORD": "pypi-flit123abc",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pypi-flit123abc",
					},
				},
			},
		},
		"HATCH_INDEX_AUTH environment variable": {
			Environment: map[string]string{
				"HATCH_INDEX_AUTH": "pypi-hatch789xyz",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pypi-hatch789xyz",
					},
				},
			},
		},
		".pypirc file with pypi section": {
			Files: map[string]string{
				"~/.pypirc": `[distutils]
index-servers = pypi

[pypi]
username = __token__
password = pypi-secret123`,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pypi-secret123",
					},
				},
			},
		},
		".pypirc file with server-login section": {
			Files: map[string]string{
				"~/.pypirc": `[server-login]
password = pypi-serverlogin456`,
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "pypi-serverlogin456",
					},
				},
			},
		},
	})
}

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default provisioner sets TWINE env vars": {
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
