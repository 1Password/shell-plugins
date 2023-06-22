package vercel

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAPITokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(
		t, APIToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
			"default": {
				ItemFields: map[sdk.FieldName]string{
					fieldname.Token: "tZk79pLyPLGgUVlkHbnLeXgl",
				},
				ExpectedOutput: sdk.ProvisionOutput{
					CommandLine: []string{"--token", "tZk79pLyPLGgUVlkHbnLeXgl"},
				},
			},
		},
	)
}

func TestAPITokenImporter(t *testing.T) {
	plugintest.TestImporter(
		t, APIToken().Importer, map[string]plugintest.ImportCase{
			"config file (macOS)": {
				OS: "darwin",
				Files: map[string]string{
					"~/Library/Application Support/com.vercel.cli/auth.json": plugintest.LoadFixture(t, "auth.json"),
				},
				ExpectedCandidates: []sdk.ImportCandidate{
					{
						Fields: map[sdk.FieldName]string{
							fieldname.Token: "tZk79pLyPLGgUVlkHbnLeXgl",
						},
					},
				},
			},
			"config file (Linux)": {
				OS: "linux",
				Files: map[string]string{
					"~/.config/com.vercel.cli/auth.json": plugintest.LoadFixture(t, "auth.json"),
				},
				ExpectedCandidates: []sdk.ImportCandidate{
					{
						Fields: map[sdk.FieldName]string{
							fieldname.Token: "tZk79pLyPLGgUVlkHbnLeXgl",
						},
					},
				},
			},
		},
	)
}
