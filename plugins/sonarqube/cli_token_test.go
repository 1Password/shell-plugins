package sonarqube

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestCLITokenImporter(t *testing.T) {
	plugintest.TestImporter(t, CLIToken().Importer, map[string]plugintest.ImportCase{
		"token only": {
			Environment: map[string]string{
				"SONARQUBE_CLI_TOKEN": "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
					},
				},
			},
		},
		"token with organization": {
			Environment: map[string]string{
				"SONARQUBE_CLI_TOKEN": "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				"SONARQUBE_CLI_ORG":   "my-org",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
						fieldname.Organization: "my-org",
					},
				},
			},
		},
		"token with server": {
			Environment: map[string]string{
				"SONARQUBE_CLI_TOKEN":  "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				"SONARQUBE_CLI_SERVER": "https://sonarqube.example.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
						fieldname.URL:   "https://sonarqube.example.com",
					},
				},
			},
		},
	})
}

func TestCLITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, CLIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"token only": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SONARQUBE_CLI_TOKEN": "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				},
			},
		},
		"token with organization": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				fieldname.Organization: "my-org",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SONARQUBE_CLI_TOKEN": "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
					"SONARQUBE_CLI_ORG":   "my-org",
				},
			},
		},
		"token with server": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				fieldname.URL:   "https://sonarqube.example.com",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"SONARQUBE_CLI_TOKEN":  "sqp_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
					"SONARQUBE_CLI_SERVER": "https://sonarqube.example.com",
				},
			},
		},
	})
}
