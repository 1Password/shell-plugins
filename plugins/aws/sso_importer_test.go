package aws

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestSSOImporter(t *testing.T) {
	plugintest.TestImporter(t, SSOProfile().Importer, map[string]plugintest.ImportCase{
		"legacy SSO profile in default config location": {
			Files: map[string]string{
				"~/.aws/config": plugintest.LoadFixture(t, "config-sso-legacy"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "sso-legacy",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:   "https://example.awsapps.com/start",
						fieldname.SSORegion:     "us-east-1",
						fieldname.SSOAccountID:  "123456789012",
						fieldname.SSORoleName:   "ReadOnly",
						fieldname.DefaultRegion: "us-west-2",
					},
				},
			},
		},
		"sso-session profile in default config location": {
			Files: map[string]string{
				"~/.aws/config": plugintest.LoadFixture(t, "config-sso-session"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "sso-session",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:   "https://corp.awsapps.com/start",
						fieldname.SSORegion:     "eu-west-1",
						fieldname.SSOAccountID:  "210987654321",
						fieldname.SSORoleName:   "Admin",
						fieldname.SSOSession:    "corp",
						fieldname.DefaultRegion: "eu-west-1",
					},
				},
			},
		},
		"mixed config: legacy + sso-session + non-SSO profile": {
			Files: map[string]string{
				"~/.aws/config": plugintest.LoadFixture(t, "config-sso-mixed"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "sso-legacy",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:   "https://example.awsapps.com/start",
						fieldname.SSORegion:     "us-east-1",
						fieldname.SSOAccountID:  "123456789012",
						fieldname.SSORoleName:   "ReadOnly",
						fieldname.DefaultRegion: "us-west-2",
					},
				},
				{
					NameHint: "sso-new",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:   "https://corp.awsapps.com/start",
						fieldname.SSORegion:     "eu-west-1",
						fieldname.SSOAccountID:  "210987654321",
						fieldname.SSORoleName:   "Admin",
						fieldname.SSOSession:    "corp",
						fieldname.DefaultRegion: "eu-west-1",
					},
				},
			},
		},
		"AWS_CONFIG_FILE env var override in home dir": {
			Environment: map[string]string{
				"AWS_CONFIG_FILE": "~/.config-custom-sso",
			},
			Files: map[string]string{
				"~/.config-custom-sso": plugintest.LoadFixture(t, "config-sso-legacy"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "sso-legacy",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:   "https://example.awsapps.com/start",
						fieldname.SSORegion:     "us-east-1",
						fieldname.SSOAccountID:  "123456789012",
						fieldname.SSORoleName:   "ReadOnly",
						fieldname.DefaultRegion: "us-west-2",
					},
				},
			},
		},
		"missing config file": {
			Files:              map[string]string{},
			ExpectedCandidates: nil,
		},
		"empty config file": {
			Files: map[string]string{
				"~/.aws/config": "",
			},
			ExpectedCandidates: nil,
		},
		"profile references unknown sso-session": {
			Files: map[string]string{
				"~/.aws/config": "[profile broken]\nsso_session = nonexistent\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"broken\" references unknown sso-session \"nonexistent\""},
							},
						},
					},
				},
			},
		},
	})
}
