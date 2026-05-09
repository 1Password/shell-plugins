package aws

import (
	"os"
	"path/filepath"
	"strings"
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
		"sso_start_url with non-HTTPS scheme is rejected": {
			Files: map[string]string{
				"~/.aws/config": "[profile bad-scheme]\nsso_start_url = http://example.awsapps.com/start\nsso_region = us-east-1\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"bad-scheme\": sso_start_url \"http://example.awsapps.com/start\" must use https://"},
							},
						},
					},
				},
			},
		},
		"sso_start_url with file:// scheme is rejected": {
			Files: map[string]string{
				"~/.aws/config": "[profile file-url]\nsso_start_url = file:///etc/passwd\nsso_region = us-east-1\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"file-url\": sso_start_url \"file:///etc/passwd\" must use https://"},
							},
						},
					},
				},
			},
		},
		"sso_account_id is not 12 digits is rejected": {
			Files: map[string]string{
				"~/.aws/config": "[profile short-acct]\nsso_start_url = https://example.awsapps.com/start\nsso_region = us-east-1\nsso_account_id = 12345\nsso_role_name = ReadOnly\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"short-acct\": sso_account_id \"12345\" is not a 12-digit AWS account ID"},
							},
						},
					},
				},
			},
		},
		"sso_region with bad characters is rejected": {
			Files: map[string]string{
				"~/.aws/config": "[profile bad-region]\nsso_start_url = https://example.awsapps.com/start\nsso_region = us@east-1\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"bad-region\": sso_region \"us@east-1\" is not a valid AWS region"},
							},
						},
					},
				},
			},
		},
		"NUL byte in sso_session is rejected": {
			Files: map[string]string{
				"~/.aws/config": "[profile nul-session]\nsso_session = corp\x00evil\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n[sso-session corp]\nsso_start_url = https://example.awsapps.com/start\nsso_region = us-east-1\n",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: sdk.ImportSource{Files: []string{"~/.aws/config"}},
						Diagnostics: sdk.Diagnostics{
							Errors: []sdk.Error{
								{Message: "profile \"nul-session\" sso_session value contains a NUL byte"},
							},
						},
					},
				},
			},
		},
		"malformed section does not brick valid profiles": {
			// `[profile valid]extra-junk` is malformed; with `Loose: true` the loader returns a
			// partial result and the well-formed profile still surfaces.
			Files: map[string]string{
				"~/.aws/config": "[profile valid]\nsso_start_url = https://example.awsapps.com/start\nsso_region = us-east-1\nsso_account_id = 123456789012\nsso_role_name = ReadOnly\n[profile valid]extra-junk\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "valid",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:  "https://example.awsapps.com/start",
						fieldname.SSORegion:    "us-east-1",
						fieldname.SSOAccountID: "123456789012",
						fieldname.SSORoleName:  "ReadOnly",
					},
				},
			},
		},
		"duplicate profile sections last-wins (botocore parity)": {
			// Last-wins matches botocore's `AllowShadows: false` default behaviour. Tested here
			// so a future loader-option change cannot silently flip the semantics.
			Files: map[string]string{
				"~/.aws/config": "[profile dup]\nsso_start_url = https://first.awsapps.com/start\nsso_region = us-east-1\nsso_account_id = 111111111111\nsso_role_name = ReadOnly\n[profile dup]\nsso_start_url = https://second.awsapps.com/start\nsso_region = eu-west-1\nsso_account_id = 222222222222\nsso_role_name = Admin\n",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					NameHint: "dup",
					Fields: map[sdk.FieldName]string{
						fieldname.SSOStartURL:  "https://second.awsapps.com/start",
						fieldname.SSORegion:    "eu-west-1",
						fieldname.SSOAccountID: "222222222222",
						fieldname.SSORoleName:  "Admin",
					},
				},
			},
		},
	})
}

// TestExternalConfigPathValidation exercises validateExternalConfigPath directly so the AWS_CONFIG_FILE
// safety properties (regular file, owned by current user, not a symlink) can be tested without
// stand-up cost in the plugintest harness.
func TestExternalConfigPathValidation(t *testing.T) {
	dir := t.TempDir()

	t.Run("regular file owned by current user is accepted", func(t *testing.T) {
		path := filepath.Join(dir, "ok-config")
		if err := os.WriteFile(path, []byte("[profile x]\n"), 0o600); err != nil {
			t.Fatalf("write fixture: %v", err)
		}
		if err := validateExternalConfigPath(path); err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("non-existent path is allowed (downstream ReadFile reports it)", func(t *testing.T) {
		if err := validateExternalConfigPath(filepath.Join(dir, "does-not-exist")); err != nil {
			t.Errorf("expected nil error for non-existent path, got %v", err)
		}
	})

	t.Run("symlink is refused", func(t *testing.T) {
		target := filepath.Join(dir, "symlink-target")
		if err := os.WriteFile(target, []byte("[profile x]\n"), 0o600); err != nil {
			t.Fatalf("write target: %v", err)
		}
		link := filepath.Join(dir, "symlink-link")
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("cannot create symlink (e.g. unprivileged Windows): %v", err)
		}
		err := validateExternalConfigPath(link)
		if err == nil {
			t.Fatal("expected symlink rejection, got nil")
		}
		if !strings.Contains(err.Error(), "symlink") {
			t.Errorf("expected error message to mention symlink, got %v", err)
		}
	})

	t.Run("directory is refused", func(t *testing.T) {
		err := validateExternalConfigPath(dir)
		if err == nil {
			t.Fatal("expected directory rejection, got nil")
		}
		if !strings.Contains(err.Error(), "regular file") {
			t.Errorf("expected error message to mention regular file, got %v", err)
		}
	})
}
