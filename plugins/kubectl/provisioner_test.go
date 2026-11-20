package kubectl

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	prodServer  = "https://prod.example.com:6443"
	localServer = "https://127.0.0.1:7443"
	itemToken   = "test-token-not-a-secret"
	testName    = "1password-shell-plugin-test"
	overlayPath = "/tmp/kubeconfig"
	fakeCertPEM = "-----BEGIN CERTIFICATE-----\nZmFrZQ==\n-----END CERTIFICATE-----\n"
	fakeKeyPEM  = "-----BEGIN PRIVATE KEY-----\nZmFrZQ==\n-----END PRIVATE KEY-----\n"
)

var (
	fakeCertB64 = base64.StdEncoding.EncodeToString([]byte(fakeCertPEM))
	fakeKeyB64  = base64.StdEncoding.EncodeToString([]byte(fakeKeyPEM))

	tokenItem = map[sdk.FieldName]string{fieldname.Address: prodServer, fieldname.Token: itemToken}
	pemItem   = map[sdk.FieldName]string{fieldname.Address: prodServer, fieldname.Certificate: fakeCertPEM, fieldname.PrivateKey: fakeKeyPEM}

	tokenUser    = "    token: " + itemToken + "\n"
	certKeyUser  = "    client-certificate-data: " + fakeCertB64 + "\n    client-key-data: " + fakeKeyB64 + "\n"
	prodCtxBody  = "    cluster: prod\n    user: " + testName + "\n"
	homeEnv      = map[string]string{"KUBECONFIG": overlayPath + ":{home}/.kube/config"}
	notTouched   = map[string]string{}
	twoContexts  = map[string]string{".kube/config": kubeconfigYAML("prod", "", "")}
	currentLocal = map[string]string{".kube/config": kubeconfigYAML("local", "", "")}
)

// kubeconfigYAML builds a config with contexts prod and local; prodUserExtra is appended to
// the prod user and prodCtxExtra to the prod context.
func kubeconfigYAML(current, prodCtxExtra, prodUserExtra string) string {
	return "apiVersion: v1\nkind: Config\n" +
		"clusters:\n" +
		"- name: prod\n  cluster:\n    server: " + prodServer + "\n" +
		"- name: local\n  cluster:\n    server: " + localServer + "\n" +
		"contexts:\n" +
		"- name: prod\n  context:\n    cluster: prod\n    user: prod-user\n" + prodCtxExtra +
		"- name: local\n  context:\n    cluster: local\n    user: local-user\n" +
		"users:\n" +
		"- name: prod-user\n  user:\n    token: stale-prod-token\n" + prodUserExtra +
		"- name: local-user\n  user:\n    token: local-token\n" +
		"current-context: " + current + "\n"
}

func overlayYAML(ctxName, ctxBody, userName, userBody string) string {
	return "apiVersion: v1\nkind: Config\n" +
		"contexts:\n- name: " + ctxName + "\n  context:\n" + ctxBody +
		"users:\n- name: " + userName + "\n  user:\n" + userBody
}

func standaloneYAML(server, caLine, userBody string) string {
	return "apiVersion: v1\nkind: Config\n" +
		"clusters:\n- name: " + testName + "\n  cluster:\n    server: " + server + "\n" + caLine +
		"contexts:\n- name: " + testName + "\n  context:\n    cluster: " + testName + "\n    user: " + testName + "\n" +
		"users:\n- name: " + testName + "\n  user:\n" + userBody +
		"current-context: " + testName + "\n"
}

func withFields(base map[sdk.FieldName]string, changes map[sdk.FieldName]string) map[sdk.FieldName]string {
	out := map[sdk.FieldName]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range changes {
		if v == "" {
			delete(out, k)
			continue
		}
		out[k] = v
	}
	return out
}

type provisionerCase struct {
	files    map[string]string // relative to the home directory; "{home}" is replaced in every string
	env      map[string]string
	args     []string
	item     map[sdk.FieldName]string
	names    []string // generator results in order; default {testName}
	wantEnv  map[string]string
	wantArgs []string // nil means args unchanged
	wantFile string   // "" means no files
	wantErrs []string
	homeDir  string // subdirectory of the test's temp dir used as home; default the temp dir itself
	tempDir  string // default "/tmp"
	random   bool   // use the default name generator instead of the injected sequence
}

func clearKubeEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"KUBECONFIG", "KUBERC", "KUBERNETES_SERVICE_HOST", "KUBERNETES_MASTER"} {
		t.Setenv(key, "")
	}
}

func runProvisionerCases(t *testing.T, cases map[string]provisionerCase) {
	t.Helper()
	clearKubeEnv(t)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			home := filepath.Join(t.TempDir(), c.homeDir)
			sub := func(s string) string { return strings.ReplaceAll(s, "{home}", home) }
			subAll := func(in []string) []string {
				out := make([]string, len(in))
				for i, s := range in {
					out[i] = sub(s)
				}
				return out
			}

			for _, key := range []string{"KUBECONFIG", "KUBERC", "KUBERNETES_SERVICE_HOST", "KUBERNETES_MASTER"} {
				t.Setenv(key, sub(c.env[key]))
			}
			for path, contents := range c.files {
				full := filepath.Join(home, path)
				require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o700))
				require.NoError(t, os.WriteFile(full, []byte(sub(contents)), 0o600))
			}

			names := c.names
			if names == nil {
				names = []string{testName}
			}
			calls := 0
			provisioner := kubeconfigProvisioner{newName: func() string {
				name := names[calls%len(names)]
				calls++
				return name
			}}
			if c.random {
				provisioner = kubeconfigProvisioner{}
			}
			tempDir := c.tempDir
			if tempDir == "" {
				tempDir = "/tmp"
			}

			args := append([]string{"kubectl"}, subAll(c.args)...)
			out := sdk.ProvisionOutput{
				Environment: map[string]string{},
				Files:       map[string]sdk.OutputFile{},
				CommandLine: append([]string(nil), args...),
			}
			provisioner.Provision(context.Background(), sdk.ProvisionInput{
				ItemFields: c.item,
				HomeDir:    home,
				TempDir:    tempDir,
			}, &out)

			want := sdk.ProvisionOutput{
				Environment: map[string]string{},
				Files:       map[string]sdk.OutputFile{},
				CommandLine: args,
			}
			for k, v := range c.wantEnv {
				want.Environment[k] = sub(v)
			}
			if c.wantArgs != nil {
				want.CommandLine = append([]string{"kubectl"}, subAll(c.wantArgs)...)
			}
			if c.wantFile != "" {
				want.Files[overlayPath] = sdk.OutputFile{Contents: []byte(c.wantFile)}
			}
			for _, msg := range c.wantErrs {
				want.Diagnostics.Errors = append(want.Diagnostics.Errors, sdk.Error{Message: msg})
			}

			assert.Equal(t, want, out)
		})
	}
}

func TestKubeconfigProvisionerMatching(t *testing.T) {
	tokenOverlay := overlayYAML("prod", prodCtxBody, testName, tokenUser)
	schemeless := map[string]string{".kube/config": strings.ReplaceAll(kubeconfigYAML("prod", "", ""), prodServer, "prod.example.com:6443")}
	schemelessItem := withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "prod.example.com:6443"})
	runProvisionerCases(t, map[string]provisionerCase{
		"P1 token item, current-context matches": {
			files: twoContexts, item: tokenItem,
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"P2 cert and key item in PEM": {
			files: twoContexts, item: pemItem,
			wantEnv: homeEnv, wantFile: overlayYAML("prod", prodCtxBody, testName, certKeyUser),
		},
		"P3 cert and key item in base64": {
			files:   twoContexts,
			item:    withFields(pemItem, map[sdk.FieldName]string{fieldname.Certificate: fakeCertB64, fieldname.PrivateKey: fakeKeyB64}),
			wantEnv: homeEnv, wantFile: overlayYAML("prod", prodCtxBody, testName, certKeyUser),
		},
		"P4 token, cert and key": {
			files: twoContexts, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.Token: itemToken}),
			wantEnv: homeEnv, wantFile: overlayYAML("prod", prodCtxBody, testName, tokenUser+certKeyUser),
		},
		"P5 current-context targets another cluster": {
			files: currentLocal, item: tokenItem, wantEnv: notTouched,
		},
		"P6 --context local": {
			files: twoContexts, item: tokenItem, args: []string{"--context", "local", "get", "pods"}, wantEnv: notTouched,
		},
		"P6 --context=local": {
			files: twoContexts, item: tokenItem, args: []string{"get", "pods", "--context=local"}, wantEnv: notTouched,
		},
		"P7 --context prod while current is local": {
			files: currentLocal, item: tokenItem, args: []string{"--context", "prod", "get", "pods"},
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"P8 -s other server": {
			files: twoContexts, item: tokenItem, args: []string{"-s", localServer, "get", "pods"}, wantEnv: notTouched,
		},
		"P8 --server=other server": {
			files: twoContexts, item: tokenItem, args: []string{"--server=" + localServer, "get", "pods"}, wantEnv: notTouched,
		},
		"P8 -s attached other server": {
			files: twoContexts, item: tokenItem, args: []string{"-s" + localServer, "get", "pods"}, wantEnv: notTouched,
		},
		"--server pointing at the item's cluster provisions the selected context": {
			files: currentLocal, item: tokenItem, args: []string{"--server", prodServer, "get", "pods"},
			wantEnv: homeEnv, wantFile: overlayYAML("local", "    cluster: local\n    user: "+testName+"\n", testName, tokenUser),
		},
		"P9 --cluster local": {
			files: twoContexts, item: tokenItem, args: []string{"--cluster", "local", "get", "pods"}, wantEnv: notTouched,
		},
		"--cluster prod while current is local keeps the original cluster in the overlay": {
			files: currentLocal, item: tokenItem, args: []string{"--cluster=prod", "get", "pods"},
			wantEnv: homeEnv, wantFile: overlayYAML("local", "    cluster: local\n    user: "+testName+"\n", testName, tokenUser),
		},
		"P10 --context after the terminator is ignored": {
			files: twoContexts, item: tokenItem, args: []string{"exec", "pod", "--", "kubectl", "--context", "local"},
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"unknown --context": {
			files: twoContexts, item: tokenItem, args: []string{"--context", "nope", "get", "pods"}, wantEnv: notTouched,
		},
		"empty current-context": {
			files: map[string]string{".kube/config": kubeconfigYAML("", "", "")}, item: tokenItem, wantEnv: notTouched,
		},
		"P19 address normalisation": {
			files:   twoContexts,
			item:    withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "HTTPS://Prod.Example.com:6443/"}),
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"http and https are different endpoints": {
			files: twoContexts, item: withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "http://prod.example.com:6443"}), wantEnv: notTouched,
		},
		"ambiguous argv provisions nothing": {
			files: twoContexts, item: tokenItem, args: []string{"get", "pods", "--template", "--server=" + prodServer}, wantEnv: notTouched,
		},
		"ambiguous short cluster provisions nothing": {
			files: twoContexts, item: tokenItem, args: []string{"get", "pods", "-As" + prodServer}, wantEnv: notTouched,
		},
		"B scheme-less kubeconfig server does not match an https Address": {
			files: schemeless, item: tokenItem, wantEnv: notTouched,
		},
		"B scheme-less kubeconfig server matches a scheme-less Address": {
			files: schemeless, item: schemelessItem,
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"B scheme-less server with :443 does not match a scheme-less Address without it": {
			files: map[string]string{".kube/config": strings.ReplaceAll(kubeconfigYAML("prod", "", ""), prodServer, "prod.example.com:443")},
			item:  withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "prod.example.com"}), wantEnv: notTouched,
		},
		"B https kubeconfig server does not match a scheme-less Address": {
			files: twoContexts, item: schemelessItem, wantEnv: notTouched,
		},
		"B scheme-less -s does not match an https Address": {
			files: twoContexts, item: tokenItem, args: []string{"-s", "prod.example.com:6443", "get", "pods"}, wantEnv: notTouched,
		},
		"B https -s matches an https Address": {
			files: currentLocal, item: tokenItem, args: []string{"-s", "https://prod.example.com:6443", "get", "pods"},
			wantEnv: homeEnv, wantFile: overlayYAML("local", "    cluster: local\n    user: "+testName+"\n", testName, tokenUser),
		},
		"F context user missing from the config": {
			files:   map[string]string{".kube/config": strings.Replace(kubeconfigYAML("prod", "", "    as: limited\n"), "user: prod-user", "user: ghost", 1)},
			item:    tokenItem,
			wantEnv: homeEnv, wantFile: overlayYAML("prod", prodCtxBody, testName, tokenUser),
		},
		"P23 context namespace copied": {
			files:   map[string]string{".kube/config": kubeconfigYAML("prod", "    namespace: team\n", "")},
			item:    tokenItem,
			wantEnv: homeEnv, wantFile: overlayYAML("prod", prodCtxBody+"    namespace: team\n", testName, tokenUser),
		},
		"P24 impersonation and user extensions copied": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "",
				"    as: limited\n    as-uid: \"42\"\n    as-groups:\n    - readers\n    as-user-extra:\n      scopes:\n      - view\n"+
					"    extensions:\n    - name: user-ext\n      extension:\n        level: 3\n")},
			item:    tokenItem,
			wantEnv: homeEnv,
			wantFile: overlayYAML("prod", prodCtxBody, testName, tokenUser+
				"    as: limited\n    as-uid: \"42\"\n    as-groups:\n    - readers\n    as-user-extra:\n      scopes:\n      - view\n"+
				"    extensions:\n    - extension:\n        level: 3\n      name: user-ext\n"),
		},
		"P25 context extensions copied": {
			files:    map[string]string{".kube/config": kubeconfigYAML("prod", "    extensions:\n    - name: ctx-ext\n      extension:\n        team: a\n", "")},
			item:     tokenItem,
			wantEnv:  homeEnv,
			wantFile: overlayYAML("prod", prodCtxBody+"    extensions:\n    - extension:\n        team: a\n      name: ctx-ext\n", testName, tokenUser),
		},
		"P26 generated name already used in the config": {
			files: map[string]string{".kube/config": strings.ReplaceAll(kubeconfigYAML("prod", "", ""), "local-user", testName)},
			item:  tokenItem, names: []string{testName, testName + "-2"},
			wantEnv: homeEnv, wantFile: overlayYAML("prod", "    cluster: prod\n    user: "+testName+"-2\n", testName+"-2", tokenUser),
		},
		"generator never yields a free name": {
			files: map[string]string{".kube/config": strings.ReplaceAll(kubeconfigYAML("prod", "", ""), "local-user", testName)},
			item:  tokenItem, wantEnv: notTouched,
			wantErrs: []string{"kubectl: could not choose a unique name for the kubeconfig overlay"},
		},
	})
}

func TestKubeconfigProvisionerExplicitAuth(t *testing.T) {
	cases := map[string]provisionerCase{}
	for name, args := range map[string][]string{
		"P11 --user":               {"--user", "x", "get", "pods"},
		"P11 --token":              {"--token", "t", "get", "pods"},
		"P11 --client-certificate": {"--client-certificate", "p", "get", "pods"},
		"--client-key":             {"--client-key=k", "get", "pods"},
		"--username":               {"--username", "u", "get", "pods"},
		"--password":               {"--password=p", "get", "pods"},
		"P32 --token= empty":       {"--token=", "get", "pods"},
		"P33 --client_key=k":       {"--client_key=k", "get", "pods"},
	} {
		cases[name] = provisionerCase{files: twoContexts, item: tokenItem, args: args, wantEnv: notTouched}
	}
	cases["explicit auth with an item without Address reports nothing"] = provisionerCase{
		files: twoContexts, item: map[sdk.FieldName]string{fieldname.Token: itemToken},
		args: []string{"--token", "t", "get", "pods"}, wantEnv: notTouched,
	}
	runProvisionerCases(t, cases)
}

func TestKubeconfigProvisionerSources(t *testing.T) {
	tokenOverlay := overlayYAML("prod", prodCtxBody, testName, tokenUser)
	prodOnly := "apiVersion: v1\nkind: Config\nclusters:\n- name: prod\n  cluster:\n    server: " + prodServer + "\n" +
		"contexts:\n- name: prod\n  context:\n    cluster: prod\n    user: prod-user\n" +
		"users:\n- name: prod-user\n  user:\n    token: stale\ncurrent-context: prod\n"
	explicitEnv := map[string]string{"KUBECONFIG": overlayPath + ":{home}/prod.yaml"}
	withProdFile := map[string]string{".kube/config": kubeconfigYAML("local", "", ""), "prod.yaml": prodOnly}

	runProvisionerCases(t, map[string]provisionerCase{
		"P12 --kubeconfig F": {
			files: withProdFile, item: tokenItem,
			args:     []string{"--kubeconfig", "{home}/prod.yaml", "get", "pods"},
			wantArgs: []string{"get", "pods"}, wantEnv: explicitEnv, wantFile: tokenOverlay,
		},
		"P12 --kubeconfig=F": {
			files: withProdFile, item: tokenItem,
			args:     []string{"get", "--kubeconfig={home}/prod.yaml", "pods", "--", "--kubeconfig", "kept"},
			wantArgs: []string{"get", "pods", "--", "--kubeconfig", "kept"}, wantEnv: explicitEnv, wantFile: tokenOverlay,
		},
		"--kubeconfig wins over KUBECONFIG": {
			files: map[string]string{"local.yaml": kubeconfigYAML("local", "", ""), "prod.yaml": prodOnly},
			env:   map[string]string{"KUBECONFIG": "{home}/local.yaml"}, item: tokenItem,
			args:     []string{"--kubeconfig", "{home}/prod.yaml", "get", "pods"},
			wantArgs: []string{"get", "pods"}, wantEnv: explicitEnv, wantFile: tokenOverlay,
		},
		"P13 --kubeconfig missing": {
			files: twoContexts, item: tokenItem, args: []string{"--kubeconfig", "{home}/missing", "get", "pods"}, wantEnv: notTouched,
		},
		"P30 --kubeconfig path with a list separator": {
			files: map[string]string{"a:b": prodOnly}, item: tokenItem, args: []string{"--kubeconfig", "{home}/a:b", "get", "pods"}, wantEnv: notTouched,
		},
		"P31 --kubeconfig directory": {
			files: twoContexts, item: tokenItem, args: []string{"--kubeconfig", "{home}/.kube", "get", "pods"}, wantEnv: notTouched,
		},
		"KUBECONFIG list": {
			files: map[string]string{"a.yaml": prodOnly, "b.yaml": kubeconfigYAML("local", "", "")},
			env:   map[string]string{"KUBECONFIG": "{home}/a.yaml:{home}/b.yaml"}, item: tokenItem,
			wantEnv: map[string]string{"KUBECONFIG": overlayPath + ":{home}/a.yaml:{home}/b.yaml"}, wantFile: tokenOverlay,
		},
		"P14 first non-empty current-context wins": {
			files: map[string]string{"a.yaml": "apiVersion: v1\nkind: Config\ncurrent-context: local\n", "b.yaml": kubeconfigYAML("prod", "", "")},
			env:   map[string]string{"KUBECONFIG": "{home}/a.yaml:{home}/b.yaml"}, item: tokenItem, wantEnv: notTouched,
		},
		"P15 first definition of a context wins": {
			files: map[string]string{
				"a.yaml": "apiVersion: v1\nkind: Config\nclusters:\n- name: elsewhere\n  cluster:\n    server: " + localServer + "\n" +
					"contexts:\n- name: prod\n  context:\n    cluster: elsewhere\n    user: prod-user\n",
				"b.yaml": kubeconfigYAML("prod", "", ""),
			},
			env: map[string]string{"KUBECONFIG": "{home}/a.yaml:{home}/b.yaml"}, item: tokenItem, wantEnv: notTouched,
		},
		"P16 empty and missing entries kept verbatim": {
			files: map[string]string{"b.yaml": kubeconfigYAML("prod", "", "")},
			env:   map[string]string{"KUBECONFIG": ":{home}/missing:{home}/b.yaml"}, item: tokenItem,
			wantEnv: map[string]string{"KUBECONFIG": overlayPath + "::{home}/missing:{home}/b.yaml"}, wantFile: tokenOverlay,
		},
		"P17 KUBECONFIG set but every file missing": {
			env: map[string]string{"KUBECONFIG": "{home}/missing-a:{home}/missing-b"}, item: tokenItem, wantEnv: notTouched,
		},
		"C home config path containing the list separator": {
			homeDir: "a:b", files: twoContexts, item: tokenItem, wantEnv: notTouched,
		},
		"F overlay path containing the list separator": {
			tempDir: "/tmp/a:b", files: twoContexts, item: tokenItem, wantEnv: notTouched,
		},
		"P18 unparseable home config": {
			files: map[string]string{".kube/config": "contexts: [not: valid"}, item: tokenItem, wantEnv: notTouched,
		},
		"KUBECONFIG with one unparseable file": {
			files: map[string]string{"a.yaml": "contexts: [not: valid", "b.yaml": kubeconfigYAML("prod", "", "")},
			env:   map[string]string{"KUBECONFIG": "{home}/a.yaml:{home}/b.yaml"}, item: tokenItem, wantEnv: notTouched,
		},
		"home config present but in-cluster variables set": {
			files: twoContexts, env: map[string]string{"KUBERNETES_SERVICE_HOST": "10.0.0.1"}, item: tokenItem,
			wantEnv: homeEnv, wantFile: tokenOverlay,
		},
	})
}

func TestKubeconfigProvisionerItemValidation(t *testing.T) {
	runProvisionerCases(t, map[string]provisionerCase{
		"P20 certificate without private key": {
			files: twoContexts, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.PrivateKey: ""}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: the 1Password item needs both Certificate and Private Key"},
		},
		"private key without certificate": {
			files: twoContexts, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.Certificate: ""}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: the 1Password item needs both Certificate and Private Key"},
		},
		"P21 invalid base64 certificate": {
			files: twoContexts, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.Certificate: "not-base64-" + itemToken}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: Certificate in the 1Password item: value is neither PEM nor base64-encoded PEM"},
		},
		"certificate in the private key field": {
			files: twoContexts, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.PrivateKey: fakeCertPEM}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: Private Key in the 1Password item: PEM block is not of type PRIVATE KEY"},
		},
		"neither token nor certificate": {
			files: twoContexts, item: map[sdk.FieldName]string{fieldname.Address: prodServer}, wantEnv: notTouched,
			wantErrs: []string{"kubectl: the 1Password item has neither Token nor Certificate and Private Key"},
		},
		"P22 no match with a broken item": {
			files: currentLocal, item: withFields(pemItem, map[sdk.FieldName]string{fieldname.PrivateKey: ""}), wantEnv: notTouched,
		},
		"item without Address": {
			files: twoContexts, item: map[sdk.FieldName]string{fieldname.Token: itemToken}, wantEnv: notTouched,
			wantErrs: []string{"kubectl: the 1Password item has no Address; set it to the cluster server URL"},
		},
		"item Address that cannot be normalised": {
			files: twoContexts, item: withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "ftp://prod.example.com:6443"}), wantEnv: notTouched,
		},
	})
}

func TestKubeconfigProvisionerKuberc(t *testing.T) {
	tokenOverlay := overlayYAML("prod", prodCtxBody, testName, tokenUser)
	runProvisionerCases(t, map[string]provisionerCase{
		"P27 kuberc alias named like the subcommand": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercAliasYAML("get")},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: notTouched,
		},
		"P28 kuberc defaults for get with server": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercDefaultsYAML("get", "server")},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: notTouched,
		},
		"P28 kuberc defaults for get with output only": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercDefaultsYAML("get", "output")},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"kuberc namespace-only default": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercDefaultsYAML("get", "namespace")},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"kuberc namespace and server defaults": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercDefaultsYAML("get", "namespace") + "  - name: server\n    default: x\n"},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: notTouched,
		},
		"P29 KUBERC=off ignores the home kuberc": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), ".kube/kuberc": kubercAliasYAML("get")},
			env:   map[string]string{"KUBERC": "off"},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: homeEnv, wantFile: tokenOverlay,
		},
		"--kuberc pointing at an interfering file": {
			files: map[string]string{".kube/config": kubeconfigYAML("prod", "", ""), "rc": kubercDefaultsYAML("get", "context")},
			item:  tokenItem, args: []string{"--kuberc", "{home}/rc", "get", "pods"}, wantEnv: notTouched,
		},
		"interfering kuberc blocks the fresh-machine branch": {
			files: map[string]string{".kube/kuberc": kubercDefaultsYAML("get", "server")},
			item:  tokenItem, args: []string{"get", "pods"}, wantEnv: notTouched,
		},
	})
}

func TestKubeconfigProvisionerFreshMachine(t *testing.T) {
	clearKubeEnv(t)

	freshEnv := map[string]string{"KUBECONFIG": overlayPath}
	tokenStandalone := standaloneYAML(prodServer, "", tokenUser)
	plugintest.TestProvisioner(t, kubeconfigProvisioner{newName: func() string { return testName }}, map[string]plugintest.ProvisionCase{
		"T1 token item": {
			ItemFields:  tokenItem,
			CommandLine: []string{"kubectl", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: freshEnv,
				CommandLine: []string{"kubectl", "get", "pods"},
				Files:       map[string]sdk.OutputFile{overlayPath: {Contents: []byte(tokenStandalone)}},
			},
		},
		"T2 certificate authority in PEM": {
			ItemFields:  withFields(tokenItem, map[sdk.FieldName]string{fieldname.CertificateAuthority: fakeCertPEM}),
			CommandLine: []string{"kubectl", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: freshEnv,
				CommandLine: []string{"kubectl", "get", "pods"},
				Files: map[string]sdk.OutputFile{overlayPath: {Contents: []byte(
					standaloneYAML(prodServer, "    certificate-authority-data: "+fakeCertB64+"\n", tokenUser))}},
			},
		},
		"T3 Address without scheme": {
			ItemFields:  withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "prod.example.com:6443"}),
			CommandLine: []string{"kubectl", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: freshEnv,
				CommandLine: []string{"kubectl", "get", "pods"},
				Files:       map[string]sdk.OutputFile{overlayPath: {Contents: []byte(tokenStandalone)}},
			},
		},
		"T4 --context given": {
			ItemFields:     tokenItem,
			CommandLine:    []string{"kubectl", "--context", "x", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{CommandLine: []string{"kubectl", "--context", "x", "get", "pods"}},
		},
		"T5 -s other server": {
			ItemFields:     tokenItem,
			CommandLine:    []string{"kubectl", "-s", localServer, "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{CommandLine: []string{"kubectl", "-s", localServer, "get", "pods"}},
		},
		"T6 item without Address": {
			ItemFields:  map[sdk.FieldName]string{fieldname.Token: itemToken},
			CommandLine: []string{"kubectl", "get", "pods"},
			ExpectedOutput: sdk.ProvisionOutput{
				CommandLine: []string{"kubectl", "get", "pods"},
				Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{{Message: "kubectl: the 1Password item has no Address; set it to the cluster server URL"}}},
			},
		},
	})
}

func TestKubeconfigProvisionerFreshMachineGuards(t *testing.T) {
	freshEnv := map[string]string{"KUBECONFIG": overlayPath}
	runProvisionerCases(t, map[string]provisionerCase{
		"T7 legacy ~/.kube/.kubeconfig exists": {
			files: map[string]string{".kube/.kubeconfig": "apiVersion: v1\nkind: Config\n"}, item: tokenItem, wantEnv: notTouched,
		},
		"T8 KUBERNETES_MASTER set": {
			env: map[string]string{"KUBERNETES_MASTER": "https://10.0.0.1"}, item: tokenItem, wantEnv: notTouched,
		},
		"KUBERNETES_SERVICE_HOST set": {
			env: map[string]string{"KUBERNETES_SERVICE_HOST": "10.0.0.1"}, item: tokenItem, wantEnv: notTouched,
		},
		"--cluster given": {
			item: tokenItem, args: []string{"--cluster", "x", "get", "pods"}, wantEnv: notTouched,
		},
		"--kubeconfig to a missing file": {
			item: tokenItem, args: []string{"--kubeconfig", "{home}/missing", "get", "pods"}, wantEnv: notTouched,
		},
		"-s matching the Address": {
			item: tokenItem, args: []string{"-s", "https://PROD.example.com:6443/", "get", "pods"},
			wantEnv: freshEnv, wantFile: standaloneYAML(prodServer, "", tokenUser),
		},
		"empty --kubeconfig= given": {
			item: tokenItem, args: []string{"--kubeconfig=", "get", "pods"}, wantEnv: notTouched,
		},
		"E Address that does not normalise": {
			item: withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "ftp://prod.example.com:6443"}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: the 1Password item's Address is not a valid https:// or http:// server URL"},
		},
		"scheme-less -s against the https server written for a scheme-less Address": {
			item: withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "prod.example.com:6443"}),
			args: []string{"-s", "prod.example.com:6443", "get", "pods"}, wantEnv: notTouched,
		},
		"https -s against the https server written for a scheme-less Address": {
			item:    withFields(tokenItem, map[sdk.FieldName]string{fieldname.Address: "prod.example.com:6443"}),
			args:    []string{"-s", "https://prod.example.com:6443", "get", "pods"},
			wantEnv: freshEnv, wantFile: standaloneYAML(prodServer, "", tokenUser),
		},
		"-s that cannot be normalised": {
			item: tokenItem, args: []string{"-s", "ftp://prod.example.com:6443", "get", "pods"}, wantEnv: notTouched,
		},
		"invalid certificate authority": {
			item: withFields(tokenItem, map[sdk.FieldName]string{fieldname.CertificateAuthority: fakeKeyPEM}), wantEnv: notTouched,
			wantErrs: []string{"kubectl: Certificate Authority in the 1Password item: PEM block is not of type CERTIFICATE"},
		},
		"cert and key item": {
			item:    pemItem,
			wantEnv: freshEnv, wantFile: standaloneYAML(prodServer, "", certKeyUser),
		},
		"item without secrets": {
			item:     map[sdk.FieldName]string{fieldname.Address: prodServer},
			wantEnv:  notTouched,
			wantErrs: []string{"kubectl: the 1Password item has neither Token nor Certificate and Private Key"},
		},
	})
}

func TestKubeconfigProvisionerDescription(t *testing.T) {
	clearKubeEnv(t)
	assert.Equal(t, "Provision a temporary kubeconfig overlay for the context that targets the item's cluster", Provisioner().Description())
}

func TestKubeconfigProvisionerRandomName(t *testing.T) {
	clearKubeEnv(t)
	name, err := kubeconfigProvisioner{}.name()
	require.NoError(t, err)
	assert.Regexp(t, `^1password-shell-plugin-[0-9a-f]{16}$`, name)
	other, err := kubeconfigProvisioner{}.name()
	require.NoError(t, err)
	assert.NotEqual(t, name, other)
}

func TestKubeconfigProvisionerNameGenerationFailure(t *testing.T) {
	original := randRead
	randRead = func([]byte) (int, error) { return 0, errors.New("entropy unavailable") }
	t.Cleanup(func() { randRead = original })

	runProvisionerCases(t, map[string]provisionerCase{
		"D overlay for a matching context": {
			files: twoContexts, item: tokenItem, random: true, wantEnv: notTouched,
			wantErrs: []string{"kubectl: could not generate a name for the kubeconfig overlay"},
		},
		"D standalone config on a fresh machine": {
			item: tokenItem, random: true, wantEnv: notTouched,
			wantErrs: []string{"kubectl: could not generate a name for the kubeconfig overlay"},
		},
	})
}
