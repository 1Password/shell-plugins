package argocd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	prodServer  = "argocd.prod.example.com"
	localServer = "localhost:8080"
	itemToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QifQ.EXAMPLE"
)

func argocdConfig(currentContext string) string {
	return "contexts:\n" +
		"- name: prod\n  server: " + prodServer + "\n" +
		"- name: local\n  server: " + localServer + "\n" +
		"current-context: " + currentContext + "\n"
}

func TestAddressAwareProvisioner(t *testing.T) {
	prodItem := map[sdk.FieldName]string{fieldname.AuthToken: itemToken, fieldname.Address: prodServer}
	provisioned := map[string]string{"ARGOCD_AUTH_TOKEN": itemToken, "ARGOCD_SERVER": prodServer}
	notProvisioned := map[string]string{}
	currentLocal := map[string]string{".config/argocd/config": argocdConfig("local")}
	currentProd := map[string]string{".config/argocd/config": argocdConfig("prod")}

	cases := map[string]struct {
		Files       map[string]string // relative to the home directory
		Env         map[string]string
		Args        []string
		ItemFields  map[sdk.FieldName]string
		ExpectedEnv map[string]string
	}{
		"provisions without an argocd config": {
			ItemFields:  prodItem,
			ExpectedEnv: provisioned,
		},
		"provisions when current-context matches the address": {
			Files:       currentProd,
			ItemFields:  prodItem,
			ExpectedEnv: provisioned,
		},
		"skips when current-context targets another server": {
			Files:       currentLocal,
			ItemFields:  prodItem,
			ExpectedEnv: notProvisioned,
		},
		"skips when --argocd-context targets another server": {
			Files:       currentProd,
			Args:        []string{"app", "list", "--argocd-context", "local"},
			ItemFields:  prodItem,
			ExpectedEnv: notProvisioned,
		},
		"provisions when --argocd-context matches the address": {
			Files:       currentLocal,
			Args:        []string{"app", "list", "--argocd-context", "prod"},
			ItemFields:  prodItem,
			ExpectedEnv: provisioned,
		},
		"skips when --server targets another server": {
			Files:       currentProd,
			Args:        []string{"app", "list", "--server=" + localServer},
			ItemFields:  prodItem,
			ExpectedEnv: notProvisioned,
		},
		"ignores --server after --": {
			Files:       currentProd,
			Args:        []string{"app", "list", "--", "--server", localServer},
			ItemFields:  prodItem,
			ExpectedEnv: provisioned,
		},
		"skips when ARGOCD_SERVER targets another server": {
			Files:       currentProd,
			Env:         map[string]string{"ARGOCD_SERVER": localServer},
			ItemFields:  prodItem,
			ExpectedEnv: notProvisioned,
		},
		"compares addresses ignoring case, scheme, trailing slash and default port": {
			Args:        []string{"app", "list", "--server", "prod.example.com:443"},
			ItemFields:  map[sdk.FieldName]string{fieldname.AuthToken: itemToken, fieldname.Address: "https://Prod.Example.com/"},
			ExpectedEnv: map[string]string{"ARGOCD_AUTH_TOKEN": itemToken, "ARGOCD_SERVER": "https://Prod.Example.com/"},
		},
		"prefers the legacy ~/.argocd config over ~/.config/argocd": {
			Files:       map[string]string{".config/argocd/config": argocdConfig("prod"), ".argocd/config": argocdConfig("local")},
			ItemFields:  prodItem,
			ExpectedEnv: notProvisioned,
		},
		"provisions when the config can't be parsed": {
			Files:       map[string]string{".config/argocd/config": "contexts: [not: valid"},
			ItemFields:  prodItem,
			ExpectedEnv: provisioned,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			for _, key := range []string{"ARGOCD_SERVER", "ARGOCD_CONFIG_DIR", "XDG_CONFIG_HOME"} {
				t.Setenv(key, c.Env[key])
			}
			for path, contents := range c.Files {
				fullPath := filepath.Join(home, path)
				require.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o700))
				require.NoError(t, os.WriteFile(fullPath, []byte(contents), 0o600))
			}

			out := sdk.ProvisionOutput{Environment: map[string]string{}, CommandLine: append([]string{"argocd"}, c.Args...)}
			AuthToken().DefaultProvisioner.Provision(context.Background(), sdk.ProvisionInput{
				HomeDir:    home,
				ItemFields: c.ItemFields,
			}, &out)

			assert.Equal(t, c.ExpectedEnv, out.Environment)
		})
	}
}
