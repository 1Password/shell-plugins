package schema

import (
	"context"
	"net/url"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAccessKey  = sdk.CredentialName("Access Key")
	testSSOProfile = sdk.CredentialName("SSO Profile")
)

type testProvisioner struct{}

func (testProvisioner) Description() string { return "Test provisioner" }

func (testProvisioner) Provision(context.Context, sdk.ProvisionInput, *sdk.ProvisionOutput) {}

func (testProvisioner) Deprovision(context.Context, sdk.DeprovisionInput, *sdk.DeprovisionOutput) {}

func testURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}

func testCredentialType(t *testing.T, name sdk.CredentialName) CredentialType {
	t.Helper()
	return CredentialType{
		Name:    name,
		DocsURL: testURL(t, "https://example.com/docs"),
		Fields: []CredentialField{
			{
				Name:                sdk.FieldName("Secret"),
				MarkdownDescription: "The secret used to authenticate to Example.",
				Secret:              true,
			},
		},
		DefaultProvisioner: testProvisioner{},
		Importer:           func(context.Context, sdk.ImportInput, *sdk.ImportOutput) {},
	}
}

// testPluginWithTwoCredentialTypes returns a plugin that offers two credential types, both of which can be
// used by its only executable.
func testPluginWithTwoCredentialTypes(t *testing.T) Plugin {
	t.Helper()

	accessKey := testCredentialType(t, testAccessKey)
	accessKey.ManagementURL = testURL(t, "https://example.com/access-keys")

	ssoProfile := testCredentialType(t, testSSOProfile)
	ssoProfile.ManagementURL = testURL(t, "https://example.com/sso")

	return Plugin{
		Name: "example",
		Platform: PlatformInfo{
			Name:     "Example",
			Homepage: testURL(t, "https://example.com"),
		},
		Credentials: []CredentialType{accessKey, ssoProfile},
		Executables: []Executable{
			{
				Name:      "Example CLI",
				Runs:      []string{"example"},
				DocsURL:   testURL(t, "https://example.com/cli"),
				NeedsAuth: func(sdk.NeedsAuthenticationInput) bool { return true },
				Uses: []CredentialUsage{
					{Name: testAccessKey, Optional: true},
					{Name: testSSOProfile, Optional: true},
				},
			},
		},
	}
}

// failedChecks returns the descriptions of all checks that did not pass, warnings included.
func failedChecks(report ValidationReport) []string {
	var failed []string
	for _, check := range report.Checks {
		if !check.Assertion {
			failed = append(failed, check.Description)
		}
	}
	return failed
}

func TestPluginValidate(t *testing.T) {
	const (
		credentialsUsedCheck   = "If more than one credential type is defined, each one is used by at least one executable"
		uniqueCredentialsCheck = "Credentials are uniquely identifiable inside a plugin"
		referencesCheck        = "Credentials referenced in executables are included in the same plugin definition"
	)

	cases := map[string]struct {
		modify         func(p *Plugin)
		expectedFailed []string
		// expectedError is whether one of the failed checks is an error instead of a warning.
		expectedError bool
	}{
		"two credential types that are both used": {
			modify: func(p *Plugin) {},
		},
		"two credential types that are used by different executables": {
			modify: func(p *Plugin) {
				other := p.Executables[0]
				other.Name = "Other Example CLI"
				other.Runs = []string{"other-example"}
				other.Uses = []CredentialUsage{{Name: testSSOProfile}}
				p.Executables[0].Uses = []CredentialUsage{{Name: testAccessKey}}
				p.Executables = append(p.Executables, other)
			},
		},
		"two credential types and no executables": {
			modify: func(p *Plugin) {
				p.Executables = nil
			},
		},
		"one credential type that is not used by the executable": {
			modify: func(p *Plugin) {
				p.Credentials = p.Credentials[:1]
				p.Executables[0].Uses = []CredentialUsage{
					{SelectFrom: &CredentialSelection{ID: "project", IncludeAllCredentials: true}},
				}
			},
		},
		"two credential types of which one is used with the plugin named explicitly": {
			modify: func(p *Plugin) {
				p.Executables[0].Uses[1].Plugin = p.Name
			},
		},
		"two credential types of which one is only reachable through a selection of all credentials": {
			modify: func(p *Plugin) {
				p.Executables[0].Uses[1] = CredentialUsage{
					SelectFrom: &CredentialSelection{ID: "project", IncludeAllCredentials: true},
				}
			},
		},
		"two credential types and an executable that only has a selection of all credentials": {
			modify: func(p *Plugin) {
				p.Executables[0].Uses = []CredentialUsage{
					{SelectFrom: &CredentialSelection{ID: "project", IncludeAllCredentials: true}},
				}
			},
		},
		"two credential types of which one is not used by any executable": {
			modify: func(p *Plugin) {
				p.Executables[0].Uses = []CredentialUsage{{Name: testAccessKey}}
			},
			expectedFailed: []string{credentialsUsedCheck},
		},
		"two credential types of which one is only used through another plugin's name": {
			modify: func(p *Plugin) {
				p.Executables[0].Uses[1].Plugin = "otherplugin"
			},
			expectedFailed: []string{credentialsUsedCheck},
		},
		"two credential types with the same name": {
			modify: func(p *Plugin) {
				p.Credentials[1].Name = testAccessKey
				p.Executables[0].Uses = []CredentialUsage{{Name: testAccessKey}}
			},
			expectedFailed: []string{uniqueCredentialsCheck},
			expectedError:  true,
		},
		"executable uses a credential type that the plugin does not define": {
			modify: func(p *Plugin) {
				p.Credentials = p.Credentials[:1]
			},
			expectedFailed: []string{referencesCheck},
			expectedError:  true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			plugin := testPluginWithTwoCredentialTypes(t)
			tc.modify(&plugin)

			valid, report := plugin.Validate()

			assert.Equal(t, tc.expectedFailed, failedChecks(report))
			assert.Equal(t, len(tc.expectedFailed) == 0, valid)
			assert.Equal(t, tc.expectedError, report.HasErrors())
		})
	}
}

func TestPluginDeepValidateWithTwoCredentialTypes(t *testing.T) {
	plugin := testPluginWithTwoCredentialTypes(t)

	reports := plugin.DeepValidate()

	// One report for the plugin, one per credential type, one for the executable and one per credential usage.
	require.Len(t, reports, 6)
	for _, report := range reports {
		assert.Empty(t, failedChecks(report), report.Heading)
	}
}

func TestPluginMarshalJSON(t *testing.T) {
	cases := map[string]struct {
		modify   func(p *Plugin)
		expected string
	}{
		"one credential type": {
			modify: func(p *Plugin) {
				p.Credentials = p.Credentials[:1]
			},
			expected: `{"platform_name":"Example","management_url":"https://example.com/access-keys"}`,
		},
		"one credential type without a management URL": {
			modify: func(p *Plugin) {
				p.Credentials = p.Credentials[:1]
				p.Credentials[0].ManagementURL = nil
			},
			expected: `{"platform_name":"Example"}`,
		},
		"two credential types use the management URL of the first": {
			modify:   func(p *Plugin) {},
			expected: `{"platform_name":"Example","management_url":"https://example.com/access-keys"}`,
		},
		"two credential types of which only the second has a management URL": {
			modify: func(p *Plugin) {
				p.Credentials[0].ManagementURL = nil
			},
			expected: `{"platform_name":"Example","management_url":"https://example.com/sso"}`,
		},
		"two credential types without a management URL": {
			modify: func(p *Plugin) {
				p.Credentials[0].ManagementURL = nil
				p.Credentials[1].ManagementURL = nil
			},
			expected: `{"platform_name":"Example"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			plugin := testPluginWithTwoCredentialTypes(t)
			tc.modify(&plugin)

			result, err := plugin.MarshalJSON()

			require.NoError(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
