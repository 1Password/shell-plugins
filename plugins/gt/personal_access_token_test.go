package gt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestHostsPath(t *testing.T) {
	t.Run("honors GITEA_CLI_CONFIG_DIR when set", func(t *testing.T) {
		giteaCliConfigDir := filepath.Join(t.TempDir(), "giteaCli")
		t.Setenv("GITEA_CLI_CONFIG_DIR", giteaCliConfigDir)

		xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
		t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

		want := filepath.Join(giteaCliConfigDir, "hosts.yml")
		if got := HostsPath(); got != want {
			t.Errorf("HostsPath() = %q, want %q", got, want)
		}
	})

	t.Run("honors XDG_CONFIG_HOME when set", func(t *testing.T) {
		// gt reads its hosts info from GITEA_CLI_CONFIG_DIR/gitea-cli/hosts.yml via
		// github.com/adrg/xdg, so the plugin must write to the same place.
		xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
		t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

		want := filepath.Join(xdgConfigHome, "gitea-cli", "hosts.yml")
		if got := HostsPath(); got != want {
			t.Errorf("HostsPath() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to os.UserConfigDir when XDG_CONFIG_HOME is unset", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")

		configDir, err := os.UserConfigDir()
		if err != nil {
			t.Skipf("os.UserConfigDir() unavailable: %v", err)
		}

		want := filepath.Join(configDir, "gitea-cli", "hosts.yml")
		if got := HostsPath(); got != want {
			t.Errorf("HostsPath() = %q, want %q", got, want)
		}
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token:       "5obcw3wk772a2beqezslabcmivghbci4qexample",
				fieldname.HostAddress: "https://git.example.com",
				fieldname.User:        "example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"GITEA_TOKEN":        "5obcw3wk772a2beqezslabcmivghbci4qexample",
					"GITEA_INSTANCE_URL": "https://git.example.com",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

	baseConfigDir := filepath.Join(xdgConfigHome, "gitea-cli")

	credFile1 := "cf196e961d99c01a4ae839aa3c7d0ba52e67aafd7ce088d5df1c96ce5eeb5dee.json"
	credFile2 := "86707cbf0f8161fb4b5ecfbaa450653003d4e0b809b06a66aac978452164d233.json"

	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"GITEA_INSTANCE_URL": {
			Environment: map[string]string{
				"GITEA_TOKEN":        "5obcw3wk772a2beqezslabcmivghbci4qexample",
				"GITEA_INSTANCE_URL": "https://git.example.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "5obcw3wk772a2beqezslabcmivghbci4qexample",
						fieldname.HostAddress: "https://git.example.com",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "5obcw3wk772a2beqezslabcmivghbci4qexample",
					},
				},
			},
		},
		"GITEA_CLI_HOST": {
			Environment: map[string]string{
				"GITEA_TOKEN":    "5obcw3wk772a2beqezslabcmivghbci4qexample",
				"GITEA_CLI_HOST": "https://git.example.com",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "5obcw3wk772a2beqezslabcmivghbci4qexample",
						fieldname.HostAddress: "https://git.example.com",
					},
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "5obcw3wk772a2beqezslabcmivghbci4qexample",
					},
				},
			},
		},
		"hosts file": {
			Files: map[string]string{
				HostsPath(): plugintest.LoadFixture(t, "hosts.yml"),
				filepath.Join(baseConfigDir, "credentials", credFile1): plugintest.LoadFixture(t, "credentials/"+credFile1),
				filepath.Join(baseConfigDir, "credentials", credFile2): plugintest.LoadFixture(t, "credentials/"+credFile2),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
						fieldname.HostAddress: "https://git.example.com",
						fieldname.User:        "example",
					},
					NameHint: "example@https://git.exa…",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "enjkarzu2ca5ffcnvzaczxncuczeoq9utlpqqrzs",
						fieldname.HostAddress: "https://gitea.com",
						fieldname.User:        "example@example.com",
					},
					NameHint: "example@example.com@htt…",
				},
			},
		},
	})
}
