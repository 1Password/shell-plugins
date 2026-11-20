package kubectl

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
)

func TestCredentialsProvisionerSmoke(t *testing.T) {
	clearKubeEnv(t)
	plugintest.TestProvisioner(t, Credentials().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"explicit context leaves the command untouched": {
			ItemFields:     tokenItem,
			CommandLine:    []string{"kubectl", "--context", "x", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{CommandLine: []string{"kubectl", "--context", "x", "get", "pods"}},
		},
	})
}

func TestCredentialsImporterSmoke(t *testing.T) {
	setImporterEnv(t, "")
	plugintest.TestImporter(t, Credentials().Importer, map[string]plugintest.ImportCase{
		"home kubeconfig": {
			Files: map[string]string{homeKubeconfig: plugintest.LoadFixture(t, "token.yaml")},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						"Address":               prodImportAddr,
						"Token":                 importToken,
						"Certificate Authority": importCAData,
					},
					NameHint: "prod",
				},
			},
		},
	})
}

func TestKubectlCLINeedsAuthSmoke(t *testing.T) {
	plugintest.TestNeedsAuth(t, KubectlCLI().NeedsAuth, map[string]plugintest.NeedsAuthCase{
		"get pods needs auth": {Args: []string{"get", "pods"}, ExpectedNeedsAuth: true},
		"help does not":       {Args: []string{"--help"}, ExpectedNeedsAuth: false},
	})
}
