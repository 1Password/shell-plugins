package npm

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "npm_example123",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"npm_config_//registry.npmjs.org/:_authToken": "npm_example123",
				},
			},
		},
		"custom default registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "custom_example123",
				fieldname.Host:  "registry.example.com/npm",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"npm_config_//registry.example.com/npm/:_authToken": "custom_example123",
					"npm_config_registry":                               "https://registry.example.com/npm/",
				},
			},
		},
		"scoped custom registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "custom_example123",
				fieldname.Host:         "https://registry.example.com/npm/",
				fieldname.Organization: "@acme",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"npm_config_//registry.example.com/npm/:_authToken": "custom_example123",
					"npm_config_@acme:registry":                         "https://registry.example.com/npm/",
				},
			},
		},
	})
}

func TestPNPMProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PNPMCLI().Uses[0].Provisioner, map[string]plugintest.ProvisionCase{
		"default registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "npm_example123",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PNPM_CONFIG__AUTH": `{"https://registry.npmjs.org/":{"@":{"authToken":"npm_example123"}}}`,
				},
			},
		},
		"custom default registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "custom_example123",
				fieldname.Host:  "registry.example.com/npm",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PNPM_CONFIG__AUTH": `{"https://registry.example.com/npm/":{"@":{"authToken":"custom_example123"}}}`,
				},
			},
		},
		"scoped custom registry": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "custom_example123",
				fieldname.Host:         "https://registry.example.com/npm/",
				fieldname.Organization: "@acme",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"PNPM_CONFIG__AUTH": `{"https://registry.example.com/npm/":{"@acme":{"authToken":"custom_example123"}}}`,
				},
			},
		},
	})
}

func TestAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AccessToken().Importer, map[string]plugintest.ImportCase{
		"NPM_TOKEN environment variable": {
			Environment: map[string]string{"NPM_TOKEN": "npm_from_env"},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{fieldname.Token: "npm_from_env"},
			}},
		},
		"NODE_AUTH_TOKEN environment variable": {
			Environment: map[string]string{"NODE_AUTH_TOKEN": "npm_from_node_env"},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{fieldname.Token: "npm_from_node_env"},
			}},
		},
		"npm user config": {
			Files: map[string]string{
				"~/.npmrc": "//registry.npmjs.org/:_authToken=npm_from_npmrc\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: "npm_from_npmrc",
					fieldname.Host:  "registry.npmjs.org",
				},
				NameHint: "registry.npmjs.org",
			}},
		},
		"scoped custom registry": {
			Files: map[string]string{
				"~/.npmrc": "@acme:registry=https://registry.example.com/npm/\n//registry.example.com/npm/:_authToken=custom_from_npmrc\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token:        "custom_from_npmrc",
					fieldname.Host:         "https://registry.example.com/npm/",
					fieldname.Organization: "acme",
				},
				NameHint: "registry.example.com/npm",
			}},
		},
		"custom default registry preserves URL": {
			Files: map[string]string{
				"~/.npmrc": "registry=http://localhost:4873/npm/\n//localhost:4873/npm/:_authToken=custom_from_npmrc\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: "custom_from_npmrc",
					fieldname.Host:  "http://localhost:4873/npm/",
				},
				NameHint: "localhost:4873/npm",
			}},
		},
		"pnpm scoped auth key": {
			OS:          "linux",
			Environment: map[string]string{"XDG_CONFIG_HOME": ""},
			Files: map[string]string{
				"~/.config/pnpm/auth.ini": "//registry.example.com/:@acme:_authToken=custom_from_pnpm\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token:        "custom_from_pnpm",
					fieldname.Host:         "registry.example.com",
					fieldname.Organization: "acme",
				},
				NameHint: "registry.example.com",
			}},
		},
		"pnpm macOS auth file": {
			OS:          "darwin",
			Environment: map[string]string{"XDG_CONFIG_HOME": ""},
			Files: map[string]string{
				"~/Library/Preferences/pnpm/auth.ini": "//registry.npmjs.org/:_authToken=npm_from_pnpm\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: "npm_from_pnpm",
					fieldname.Host:  "registry.npmjs.org",
				},
				NameHint: "registry.npmjs.org",
			}},
		},
		"pnpm XDG auth file": {
			OS:          "linux",
			Environment: map[string]string{"XDG_CONFIG_HOME": "/custom-config"},
			Files: map[string]string{
				"/custom-config/pnpm/auth.ini": "//registry.npmjs.org/:_authToken=npm_from_xdg\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: "npm_from_xdg",
					fieldname.Host:  "registry.npmjs.org",
				},
				NameHint: "registry.npmjs.org",
			}},
		},
		"legacy pnpm config": {
			Files: map[string]string{
				"~/.config/pnpm/rc": "//registry.npmjs.org/:_authToken=npm_from_legacy_pnpm\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{{
				Fields: map[sdk.FieldName]string{
					fieldname.Token: "npm_from_legacy_pnpm",
					fieldname.Host:  "registry.npmjs.org",
				},
				NameHint: "registry.npmjs.org",
			}},
		},
		"environment variable reference is not a secret": {
			Files: map[string]string{
				"~/.npmrc": "//registry.npmjs.org/:_authToken=${NPM_TOKEN}\n",
			},
			ExpectedCandidates: nil,
		},
	})
}
