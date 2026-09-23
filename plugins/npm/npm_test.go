package npm

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestNPMCLINeedsAuth(t *testing.T) {
	testPackageManagerNeedsAuth(t, NPMCLI().NeedsAuth)
}

func TestPNPMCLINeedsAuth(t *testing.T) {
	testPackageManagerNeedsAuth(t, PNPMCLI().NeedsAuth)
}

func testPackageManagerNeedsAuth(t *testing.T, needsAuth sdk.NeedsAuthentication) {
	t.Helper()

	var needsAuthCases = map[string]plugintest.NeedsAuthCase{
		"install may access private packages": {Args: []string{"install"}, ExpectedNeedsAuth: true},
		"install alias may access private packages": {
			Args:              []string{"i"},
			ExpectedNeedsAuth: true,
		},
		"publish requires authentication":  {Args: []string{"publish"}, ExpectedNeedsAuth: true},
		"whoami requires authentication":   {Args: []string{"whoami"}, ExpectedNeedsAuth: true},
		"login establishes authentication": {Args: []string{"login"}, ExpectedNeedsAuth: false},
		"run is local":                     {Args: []string{"run", "test"}, ExpectedNeedsAuth: false},
		"help does not require auth":       {Args: []string{"install", "--help"}, ExpectedNeedsAuth: false},
		"version does not require auth":    {Args: []string{"--version"}, ExpectedNeedsAuth: false},
	}

	plugintest.TestNeedsAuth(t, needsAuth, needsAuthCases)
}
