package netlify

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token: "kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6oP0qR4sT8uVexample",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"NETLIFY_AUTH_TOKEN": "kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6oP0qR4sT8uVexample",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{
				"NETLIFY_AUTH_TOKEN": "kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6oP0qR4sT8uV",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6oP0qR4sT8uV",
					},
				},
			},
		},
		"config file (macOS)": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Preferences/netlify/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "nfp_kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6EXAMPLE",
					},
				},
			},
		},
		"config file (Linux)": {
			OS: "linux",
			Files: map[string]string{
				"~/.config/netlify/config.json": plugintest.LoadFixture(t, "config.json"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "nfp_kV9wTq3ZtN0aB4cD7eF1gH5iJ8kL2mN6EXAMPLE",
					},
				},
			},
		},
		"config file (logged out)": {
			OS: "darwin",
			Files: map[string]string{
				"~/Library/Preferences/netlify/config.json": `{"telemetryDisabled": false, "cliId": "55f2daf5-ee38-4fa0-8fd8-e67892a123d6", "userId": null}`,
			},
			ExpectedCandidates: nil,
		},
	})
}
