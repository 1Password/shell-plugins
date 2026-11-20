package kubectl

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kubercHeader = "apiVersion: kubectl.config.k8s.io/v1beta1\nkind: Preference\n"

func kubercDefaultsYAML(command, option string) string {
	return kubercHeader + "defaults:\n- command: " + command + "\n  options:\n  - name: " + option + "\n    default: x\n"
}

func kubercAliasYAML(name string) string {
	return kubercHeader + "aliases:\n- name: " + name + "\n  command: get\n  prependArgs:\n  - pods\n"
}

func TestKubercInterferes(t *testing.T) {
	clearKubeEnv(t)
	const homeKuberc = ".kube/kuberc"
	cases := map[string]struct {
		files map[string]string
		fifo  string
		env   string
		args  []string
		want  bool
	}{
		"no kuberc anywhere": {
			args: []string{"get", "pods"},
		},
		"alias named like the subcommand": {
			files: map[string]string{homeKuberc: kubercAliasYAML("get")},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"alias named like a custom subcommand": {
			files: map[string]string{homeKuberc: kubercAliasYAML("getp")},
			args:  []string{"--context", "prod", "getp"},
			want:  true,
		},
		"alias named like a later positional": {
			files: map[string]string{homeKuberc: kubercAliasYAML("pods")},
			args:  []string{"-o", "yaml", "get", "pods"},
			want:  true,
		},
		"alias named like a word after the terminator": {
			files: map[string]string{homeKuberc: kubercAliasYAML("sh")},
			args:  []string{"exec", "pod", "--", "sh"},
		},
		"alias with another name": {
			files: map[string]string{homeKuberc: kubercAliasYAML("getn")},
			args:  []string{"get", "pods"},
		},
		"defaults with a harmless option": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "output")},
			args:  []string{"get", "pods"},
		},
		"defaults for another subcommand with server": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "server")},
			args:  []string{"describe", "pods"},
			want:  true,
		},
		"get server default with a local flag value before the subcommand": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "server")},
			args:  []string{"-o", "yaml", "get", "pods"},
			want:  true,
		},
		"create configmap default reached through the cm alias": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("create configmap", "context")},
			args:  []string{"create", "cm", "x"},
			want:  true,
		},
		"defaults for a two-word command": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("rollout status", "context")},
			args:  []string{"rollout", "history", "deploy/x"},
			want:  true,
		},
		"namespace-only default": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "namespace")},
			args:  []string{"get", "pods"},
		},
		"namespace shorthand default": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "n")},
			args:  []string{"get", "pods"},
		},
		"namespace and server defaults": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "namespace") + "  - name: server\n    default: x\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"option name with underscore": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "client_key")},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"option name as shorthand": {
			files: map[string]string{homeKuberc: kubercDefaultsYAML("get", "s")},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"v1alpha1 overrides with flags": {
			files: map[string]string{homeKuberc: "apiVersion: kubectl.config.k8s.io/v1alpha1\nkind: Preference\noverrides:\n- command: get\n  flags:\n  - name: server\n    default: https://127.0.0.1:7443\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"v1alpha1 overrides with a harmless flag": {
			files: map[string]string{homeKuberc: "apiVersion: kubectl.config.k8s.io/v1alpha1\nkind: Preference\noverrides:\n- command: get\n  flags:\n  - name: output\n    default: wide\n"},
			args:  []string{"get", "pods"},
		},
		"interfering second document": {
			files: map[string]string{homeKuberc: kubercHeader + "---\n" + kubercDefaultsYAML("get", "context")},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"wrong apiVersion": {
			files: map[string]string{homeKuberc: "apiVersion: v1\nkind: Preference\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"wrong kind": {
			files: map[string]string{homeKuberc: "apiVersion: kubectl.config.k8s.io/v1beta1\nkind: Config\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"second document without apiVersion": {
			files: map[string]string{homeKuberc: kubercHeader + "---\naliases: []\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"empty file": {
			files: map[string]string{homeKuberc: ""},
			args:  []string{"get", "pods"},
		},
		"unparseable file": {
			files: map[string]string{homeKuberc: "aliases: [not: valid"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"wrong shape": {
			files: map[string]string{homeKuberc: kubercHeader + "defaults: just-a-string\n"},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"KUBERC=off ignores the home kuberc": {
			files: map[string]string{homeKuberc: kubercAliasYAML("get")},
			env:   "off",
			args:  []string{"get", "pods"},
		},
		"KUBERC points at an interfering file": {
			files: map[string]string{"custom/kuberc": kubercDefaultsYAML("get", "context")},
			env:   "{home}/custom/kuberc",
			args:  []string{"get", "pods"},
			want:  true,
		},
		"KUBERC points at a missing file and the home kuberc is not read": {
			files: map[string]string{homeKuberc: kubercAliasYAML("get")},
			env:   "{home}/missing",
			args:  []string{"get", "pods"},
		},
		"--kuberc wins over KUBERC": {
			files: map[string]string{"benign": kubercHeader, "bad": kubercDefaultsYAML("get", "kubeconfig")},
			env:   "{home}/benign",
			args:  []string{"--kuberc", "{home}/bad", "get", "pods"},
			want:  true,
		},
		"--kuberc pointing at a benign file shadows the home kuberc": {
			files: map[string]string{homeKuberc: kubercAliasYAML("get"), "benign": kubercHeader},
			args:  []string{"--kuberc={home}/benign", "get", "pods"},
		},
		"--kuberc given twice": {
			files: map[string]string{"benign": kubercHeader},
			args:  []string{"--kuberc", "{home}/benign", "--kuberc={home}/benign", "get", "pods"},
			want:  true,
		},
		"kuberc is a directory": {
			files: map[string]string{homeKuberc + "/placeholder": ""},
			args:  []string{"get", "pods"},
			want:  true,
		},
		"kuberc is a FIFO": {
			fifo: "fifo",
			env:  "{home}/fifo",
			args: []string{"get", "pods"},
			want: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("KUBERC", strings.ReplaceAll(c.env, "{home}", home))
			for path, contents := range c.files {
				full := filepath.Join(home, path)
				require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o700))
				require.NoError(t, os.WriteFile(full, []byte(contents), 0o600))
			}
			if c.fifo != "" {
				if err := exec.Command("mkfifo", filepath.Join(home, c.fifo)).Run(); err != nil {
					t.Skipf("mkfifo unavailable: %v", err)
				}
			}
			args := make([]string, len(c.args))
			for i, a := range c.args {
				args[i] = strings.ReplaceAll(a, "{home}", home)
			}

			assert.Equal(t, c.want, kubercInterferes(home, parseArgs(args)))
		})
	}
}

func TestKubercInterferesForEveryTargetingOption(t *testing.T) {
	clearKubeEnv(t)
	for _, option := range []string{
		"context", "cluster", "server", "kubeconfig", "user", "token", "client-certificate", "client-key",
		"username", "password", "proxy-url", "tls-server-name", "insecure-skip-tls-verify", "certificate-authority",
	} {
		t.Run(option, func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("KUBERC", "")
			require.NoError(t, os.MkdirAll(filepath.Join(home, ".kube"), 0o700))
			require.NoError(t, os.WriteFile(filepath.Join(home, ".kube", "kuberc"), []byte(kubercDefaultsYAML("apply", option)), 0o600))

			assert.True(t, kubercInterferes(home, parseArgs([]string{"get", "pods"})))
		})
	}
}
