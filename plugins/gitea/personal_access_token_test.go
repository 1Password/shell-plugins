package gitea

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestConfigPath(t *testing.T) {
	t.Run("honors XDG_CONFIG_HOME when set", func(t *testing.T) {
		// tea reads its config from $XDG_CONFIG_HOME/tea/config.yml via
		// github.com/adrg/xdg, so the plugin must write to the same place.
		xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
		t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

		want := filepath.Join(xdgConfigHome, "tea", "config.yml")
		if got := ConfigPath(); got != want {
			t.Errorf("ConfigPath() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to os.UserConfigDir when XDG_CONFIG_HOME is unset", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")

		configDir, err := os.UserConfigDir()
		if err != nil {
			t.Skipf("os.UserConfigDir() unavailable: %v", err)
		}

		want := filepath.Join(configDir, "tea", "config.yml")
		if got := ConfigPath(); got != want {
			t.Errorf("ConfigPath() = %q, want %q", got, want)
		}
	})
}

func TestPersonalAccessTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
				fieldname.HostAddress: "https://git.example.com",
				fieldname.User:        "example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					ConfigPath(): {
						Contents: []byte(plugintest.LoadFixture(t, "config.yml")),
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenProvisionerHonorsXDGConfigHome(t *testing.T) {
	xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

	plugintest.TestProvisioner(t, PersonalAccessToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"xdg config path": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
				fieldname.HostAddress: "https://git.example.com",
				fieldname.User:        "example",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Files: map[string]sdk.OutputFile{
					filepath.Join(xdgConfigHome, "tea", "config.yml"): {
						Contents: []byte(plugintest.LoadFixture(t, "config.yml")),
					},
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{

		"config file": {
			Files: map[string]string{
				ConfigPath(): plugintest.LoadFixture(t, "import_config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
						fieldname.HostAddress: "https://git.example.com",
						fieldname.User:        "example",
					},
					NameHint: "git.example.com",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "enjkarzu2ca5ffcnvzaczxncuczeoq9utlpqqrzs",
						fieldname.HostAddress: "https://gitea.com",
						fieldname.User:        "example@example.com",
					},
					NameHint: "gitea.com",
				},
			},
		},
	})
}

func TestPersonalAccessTokenImporterHonorsXDGConfigHome(t *testing.T) {
	xdgConfigHome := filepath.Join(t.TempDir(), "xdg")
	t.Setenv("XDG_CONFIG_HOME", xdgConfigHome)

	plugintest.TestImporter(t, PersonalAccessToken().Importer, map[string]plugintest.ImportCase{
		"xdg config file": {
			Files: map[string]string{
				filepath.Join(xdgConfigHome, "tea", "config.yml"): plugintest.LoadFixture(t, "import_config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "oyyfsny27bgphldmhvffxhhlmqvdkzjrfslrsj9f",
						fieldname.HostAddress: "https://git.example.com",
						fieldname.User:        "example",
					},
					NameHint: "git.example.com",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:       "enjkarzu2ca5ffcnvzaczxncuczeoq9utlpqqrzs",
						fieldname.HostAddress: "https://gitea.com",
						fieldname.User:        "example@example.com",
					},
					NameHint: "gitea.com",
				},
			},
		},
	})
}
