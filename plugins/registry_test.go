package plugins

import (
	"net/url"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePlugins(t *testing.T) {
	for _, p := range registry {
		_, report := p.Validate()
		for _, c := range report.Checks {
			if !c.Assertion && c.Severity == schema.ValidationSeverityError {
				t.Logf("The '%s' plugin has validation errors: %s", p.Name, c.Description)
				t.Fail()
			}
		}
	}
}

func TestAllPluginsHaveUniqueNames(t *testing.T) {
	var pluginNames []string
	for _, p := range registry {
		pluginNames = append(pluginNames, p.Name)
	}

	assert.True(t, schema.IsStringSliceASet(pluginNames))
}

// useTestRegistry replaces the registered plugins with the given ones for the duration of the test.
func useTestRegistry(t *testing.T, testPlugins ...schema.Plugin) {
	t.Helper()
	original := registry
	registry = testPlugins
	t.Cleanup(func() {
		registry = original
	})
}

func testURL(t *testing.T, rawURL string) *url.URL {
	t.Helper()
	u, err := url.Parse(rawURL)
	require.NoError(t, err)
	return u
}

func testPlugins(t *testing.T) []schema.Plugin {
	t.Helper()
	return []schema.Plugin{
		{
			Name:     "single",
			Platform: schema.PlatformInfo{Name: "Single"},
			Credentials: []schema.CredentialType{
				{Name: sdk.CredentialName("API Key"), ManagementURL: testURL(t, "https://single.example.com/keys")},
			},
		},
		{
			Name:     "double",
			Platform: schema.PlatformInfo{Name: "Double"},
			Credentials: []schema.CredentialType{
				{Name: sdk.CredentialName("Access Key")},
				{Name: sdk.CredentialName("SSO Profile"), ManagementURL: testURL(t, "https://double.example.com/sso")},
			},
		},
		{
			Name:     "nocredentials",
			Platform: schema.PlatformInfo{Name: "No Credentials"},
		},
	}
}

func TestGetCredentialType(t *testing.T) {
	useTestRegistry(t, testPlugins(t)...)

	cases := map[string]struct {
		pluginName     string
		credentialName string
		expectedName   sdk.CredentialName
		expectedErr    string
	}{
		"only credential type of a plugin": {
			pluginName:     "single",
			credentialName: "API Key",
			expectedName:   "API Key",
		},
		"first credential type of a plugin with two": {
			pluginName:     "double",
			credentialName: "Access Key",
			expectedName:   "Access Key",
		},
		"second credential type of a plugin with two": {
			pluginName:     "double",
			credentialName: "SSO Profile",
			expectedName:   "SSO Profile",
		},
		"credential type identified by its ID": {
			pluginName:     "double",
			credentialName: "sso_profile",
			expectedName:   "SSO Profile",
		},
		"credential type the plugin does not have": {
			pluginName:     "single",
			credentialName: "Access Key",
			expectedErr:    "unknown credential: Access Key (single)",
		},
		"plugin without credential types": {
			pluginName:     "nocredentials",
			credentialName: "Access Key",
			expectedErr:    "unknown credential: Access Key (nocredentials)",
		},
		"plugin that is not registered": {
			pluginName:     "missing",
			credentialName: "Access Key",
			expectedErr:    "unknown plugin: missing",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			credential, err := GetCredentialType(tc.pluginName, tc.credentialName)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedName, credential.Name)
		})
	}
}

func TestRegistryJSON(t *testing.T) {
	useTestRegistry(t, testPlugins(t)...)

	result, err := RegistryJSON()

	require.NoError(t, err)
	// A plugin gets a single entry, no matter how many credential types it has, with the management URL
	// of the first credential type that has one. Plugins without credential types are left out.
	assert.JSONEq(t, `[
		{"platform_name": "Single", "management_url": "https://single.example.com/keys"},
		{"platform_name": "Double", "management_url": "https://double.example.com/sso"}
	]`, string(result))
}
