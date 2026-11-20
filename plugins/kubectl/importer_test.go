package kubectl

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const (
	importToken    = "test-token-not-a-secret"
	importCertData = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tClptRnJaUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
	importKeyData  = "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tClptRnJaUT09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K"
	importCAData   = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tClptRnJaUzFqWVE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
	prodImportAddr = "https://prod.example.com:6443"
	devImportAddr  = "https://dev.example.com:6443"
	homeKubeconfig = "~/.kube/config"
)

func setImporterEnv(t *testing.T, kubeconfigValue string) {
	t.Setenv("KUBECONFIG", kubeconfigValue)
	t.Setenv("KUBERC", "")
	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	t.Setenv("KUBERNETES_MASTER", "")
}

func kubeconfigList(paths ...string) string {
	return strings.Join(paths, string(filepath.ListSeparator))
}

func fixtureBase64(t *testing.T, name string) string {
	return base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, name)))
}

func importErrors(messages ...string) sdk.Diagnostics {
	errs := make([]sdk.Error, 0, len(messages))
	for _, m := range messages {
		errs = append(errs, sdk.Error{Message: m})
	}
	return sdk.Diagnostics{Errors: errs}
}

func tokenCandidate(address, token, nameHint string) sdk.ImportCandidate {
	return sdk.ImportCandidate{
		Fields: map[sdk.FieldName]string{
			fieldname.Address: address,
			fieldname.Token:   token,
		},
		NameHint: nameHint,
	}
}

func TestImporterHomeConfig(t *testing.T) {
	setImporterEnv(t, "")

	plugintest.TestImporter(t, TryKubeconfigFiles(), map[string]plugintest.ImportCase{
		"I1 token user with certificate authority data": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address:              prodImportAddr,
						fieldname.Token:                importToken,
						fieldname.CertificateAuthority: importCAData,
					},
					NameHint: "prod",
				},
			},
		},
		"I2 client certificate and key data stored verbatim": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "cert-data.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address:     devImportAddr,
						fieldname.Certificate: importCertData,
						fieldname.PrivateKey:  importKeyData,
					},
					NameHint: "dev",
				},
			},
		},
		"I3 relative and absolute file references are read and base64 encoded": {
			Files: map[string]string{
				homeKubeconfig:             plugintest.LoadFixture(t, "cert-paths.yaml"),
				"~/.kube/certs/client.crt": plugintest.LoadFixture(t, "client.crt"),
				"~/.kube/certs/ca.crt":     plugintest.LoadFixture(t, "ca.crt"),
				"/etc/k8s/client.key":      plugintest.LoadFixture(t, "client.key"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address:              devImportAddr,
						fieldname.Certificate:          base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, "client.crt"))),
						fieldname.PrivateKey:           base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, "client.key"))),
						fieldname.CertificateAuthority: base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, "ca.crt"))),
					},
					NameHint: "dev",
				},
			},
		},
		"I4 exec, auth-provider, basic auth, tokenFile and empty users are skipped": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "unsupported.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
		"I5 mixed file yields only the static token context": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "mixed.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: prodImportAddr,
						fieldname.Token:   importToken,
					},
					NameHint: "prod",
				},
			},
		},
		"I8 contexts sharing user and cluster give one candidate": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "shared-user.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: prodImportAddr,
						fieldname.Token:   importToken,
					},
					NameHint: "prod",
				},
			},
		},
		"I9 insecure cluster gives a candidate without certificate authority": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "insecure.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: "https://127.0.0.1:7443",
						fieldname.Token:   importToken,
					},
					NameHint: "local",
				},
			},
		},
		"I10 context named default gets an empty name hint": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "default-context.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: prodImportAddr,
						fieldname.Token:   importToken,
					},
					NameHint: "",
				},
			},
		},
		"I15 user with token and tokenFile is skipped": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token-and-token-file.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{},
		},
		"I16 token with a certificate but no key keeps only the token": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token-cert-only.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address: prodImportAddr,
						fieldname.Token:   importToken,
					},
					NameHint: "prod",
				},
			},
		},
		"inline PEM and wrapped base64 data are stored as canonical base64": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "cert-data-pem.yaml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Address:     devImportAddr,
						fieldname.Certificate: importCertData,
						fieldname.PrivateKey:  importKeyData,
					},
					NameHint: "dev",
				},
			},
		},
		"no kubeconfig at all": {
			ExpectedOutput: &sdk.ImportOutput{},
		},
	})
}

func TestImporterKubeconfigEnv(t *testing.T) {
	setImporterEnv(t, "")

	plugintest.TestImporter(t, TryKubeconfigFiles(), map[string]plugintest.ImportCase{
		"I6 KUBECONFIG files are merged and the home config is not read": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("", "/kube/a.yaml", "/kube/missing.yaml", "", "/kube/b.yaml", ""),
			},
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token.yaml"),
				"/kube/a.yaml": plugintest.LoadFixture(t, "merge-a.yaml"),
				"/kube/b.yaml": plugintest.LoadFixture(t, "merge-b.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address: "https://alpha.example.com:6443",
									fieldname.Token:   "test-token-alpha",
								},
								NameHint: "alpha",
							},
						},
					},
					{
						Source: importer.SourceFile("/kube/b.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address:     "https://beta.example.com:6443",
									fieldname.Certificate: importCertData,
									fieldname.PrivateKey:  importKeyData,
								},
								NameHint: "beta",
							},
						},
					},
				},
			},
		},
		"I7 KUBECONFIG entry equal to the home path appears once": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/~/.kube/config", "/~/.kube/config"),
			},
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/~/.kube/config"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address:              prodImportAddr,
									fieldname.Token:                importToken,
									fieldname.CertificateAuthority: importCAData,
								},
								NameHint: "prod",
							},
						},
					},
				},
			},
		},
		"KUBECONFIG with only missing files yields nothing": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/missing.yaml"),
			},
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{},
		},
		"I13 context resolves user and cluster defined in another file": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a.yaml", "/kube/b.yaml"),
			},
			Files: map[string]string{
				"/kube/a.yaml": plugintest.LoadFixture(t, "cross-a.yaml"),
				"/kube/b.yaml": plugintest.LoadFixture(t, "cross-b.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address: "https://remote.example.com:6443",
									fieldname.Token:   importToken,
								},
								NameHint: "cross",
							},
						},
					},
					{
						Source: importer.SourceFile("/kube/b.yaml"),
					},
				},
			},
		},
		"I14 first cluster definition wins and the candidate goes to the context's file": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a.yaml", "/kube/b.yaml"),
			},
			Files: map[string]string{
				"/kube/a.yaml": plugintest.LoadFixture(t, "conflict-a.yaml"),
				"/kube/b.yaml": plugintest.LoadFixture(t, "conflict-b.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a.yaml"),
					},
					{
						Source: importer.SourceFile("/kube/b.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address: prodImportAddr,
									fieldname.Token:   importToken,
								},
								NameHint: "shared",
							},
						},
					},
				},
			},
		},
		"cross-file relative references resolve against the file defining the user or cluster": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a/a.yaml", "/kube/b/b.yaml"),
			},
			Files: map[string]string{
				"/kube/a/a.yaml":          plugintest.LoadFixture(t, "cross-ref-a.yaml"),
				"/kube/b/b.yaml":          plugintest.LoadFixture(t, "cross-ref-b.yaml"),
				"/kube/b/ca.crt":          plugintest.LoadFixture(t, "ca.crt"),
				"/kube/b/client.crt":      plugintest.LoadFixture(t, "client.crt"),
				"/kube/b/keys/client.key": plugintest.LoadFixture(t, "client.key"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a/a.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address:              "https://remote.example.com:6443",
									fieldname.Certificate:          fixtureBase64(t, "client.crt"),
									fieldname.PrivateKey:           fixtureBase64(t, "client.key"),
									fieldname.CertificateAuthority: fixtureBase64(t, "ca.crt"),
								},
								NameHint: "cross-ref",
							},
						},
					},
					{
						Source: importer.SourceFile("/kube/b/b.yaml"),
					},
				},
			},
		},
		"first user definition wins": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a.yaml", "/kube/b.yaml"),
			},
			Files: map[string]string{
				"/kube/a.yaml": plugintest.LoadFixture(t, "user-conflict-a.yaml"),
				"/kube/b.yaml": plugintest.LoadFixture(t, "user-conflict-b.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a.yaml"),
					},
					{
						Source:     importer.SourceFile("/kube/b.yaml"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, "test-token-first", "prod")},
					},
				},
			},
		},
		"first context definition wins": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a.yaml", "/kube/b.yaml"),
			},
			Files: map[string]string{
				"/kube/a.yaml": plugintest.LoadFixture(t, "context-conflict-a.yaml"),
				"/kube/b.yaml": plugintest.LoadFixture(t, "context-conflict-b.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile("/kube/a.yaml"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "shared")},
					},
					{
						Source: importer.SourceFile("/kube/b.yaml"),
					},
				},
			},
		},
		"a KUBECONFIG entry that is a directory aborts the run": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/dir", "", "/kube/a.yaml"),
			},
			Files: map[string]string{
				"/kube/dir/placeholder": "x",
				"/kube/a.yaml":          plugintest.LoadFixture(t, "merge-a.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:      importer.SourceFile("/kube/dir"),
						Diagnostics: importErrors("reading /kube/dir: not a readable kubeconfig"),
					},
					{
						Source: importer.SourceFile("/kube/a.yaml"),
					},
				},
			},
		},
		"I12 a parse error in any file yields no candidates": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("/kube/a.yaml", "/kube/broken.yaml"),
			},
			Files: map[string]string{
				"/kube/a.yaml":      plugintest.LoadFixture(t, "merge-a.yaml"),
				"/kube/broken.yaml": plugintest.LoadFixture(t, "invalid.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("/kube/a.yaml"),
					},
					{
						Source: importer.SourceFile("/kube/broken.yaml"),
						Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{
							{Message: "parsing /kube/broken.yaml: not a valid kubeconfig"},
						}},
					},
				},
			},
		},
	})
}

func TestImporterDiagnostics(t *testing.T) {
	setImporterEnv(t, "")

	plugintest.TestImporter(t, TryKubeconfigFiles(), map[string]plugintest.ImportCase{
		"I11 missing referenced key file skips the context with an error": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "missing-key-file.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile(homeKubeconfig),
						Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{
							{Message: `context "dev": client key at keys/missing.key could not be imported`},
						}},
					},
				},
			},
		},
		"I12 invalid YAML": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "invalid.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile(homeKubeconfig),
						Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{
							{Message: "parsing ~/.kube/config: not a valid kubeconfig"},
						}},
					},
				},
			},
		},
		"I19 unreadable or invalid certificate authority is omitted with an error naming only the location": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "unreadable-ca.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile(homeKubeconfig),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address: prodImportAddr,
									fieldname.Token:   importToken,
								},
								NameHint: "prod",
							},
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address: devImportAddr,
									fieldname.Token:   importToken,
								},
								NameHint: "dev",
							},
						},
						Diagnostics: sdk.Diagnostics{Errors: []sdk.Error{
							{Message: `context "prod": certificate authority at missing-ca.crt could not be imported`},
							{Message: `context "dev": certificate-authority-data in ~/.kube/config could not be imported`},
						}},
					},
				},
			},
		},
		"I17 invalid key data drops the pair, keeps the token and never echoes the value": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "invalid-key-data.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile(homeKubeconfig),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "prod")},
						Diagnostics: importErrors(
							`context "prod": client key in ~/.kube/config could not be imported`,
							`context "dev": client key in ~/.kube/config could not be imported`,
						),
					},
				},
			},
		},
		"unreadable key file with a token keeps the token and reports once": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "missing-key-with-token.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:      importer.SourceFile(homeKubeconfig),
						Candidates:  []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "prod")},
						Diagnostics: importErrors(`context "prod": client key at keys/missing.key could not be imported`),
					},
				},
			},
		},
		"a context dropped as a duplicate reports nothing": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "duplicate-with-diagnostic.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile(homeKubeconfig),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "plain")},
					},
				},
			},
		},
		"a broken certificate authority shared by two candidates is reported once": {
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "shared-broken-ca.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile(homeKubeconfig),
						Candidates: []sdk.ImportCandidate{
							tokenCandidate(prodImportAddr, "test-token-one", "one"),
							tokenCandidate(prodImportAddr, "test-token-two", "two"),
						},
						Diagnostics: importErrors(`context "one": certificate authority at missing-ca.crt could not be imported`),
					},
				},
			},
		},
		"invalid file content is reported and secret-looking paths are not echoed": {
			Files: map[string]string{
				homeKubeconfig:         plugintest.LoadFixture(t, "invalid-material.yaml"),
				"~/.kube/bad-ca.crt":   plugintest.LoadFixture(t, "not-pem.txt"),
				"~/.kube/bad-cert.crt": plugintest.LoadFixture(t, "not-pem.txt"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile(homeKubeconfig),
						Candidates: []sdk.ImportCandidate{
							tokenCandidate(prodImportAddr, "test-token-one", "pem-path"),
							tokenCandidate(devImportAddr, "test-token-two", "bad-file"),
						},
						Diagnostics: importErrors(
							`context "pem-path": client key at <invalid path> could not be imported`,
							`context "pem-path": certificate authority at <invalid path> could not be imported`,
							`context "bad-file": client certificate in bad-cert.crt could not be imported`,
							`context "bad-file": certificate authority in bad-ca.crt could not be imported`,
						),
					},
				},
			},
		},
	})
}

func TestImportDisplayPath(t *testing.T) {
	setImporterEnv(t, "")

	for in, want := range map[string]string{
		"certs/client.crt":                 "certs/client.crt",
		"/etc/k8s/client.key":              "/etc/k8s/client.key",
		"certs/a\nb.crt":                   "<invalid path>",
		"-----BEGIN PRIVATE KEY-----":      "<invalid path>",
		"LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0t": "<invalid path>",
		strings.Repeat("a", 512):           strings.Repeat("a", 512),
		strings.Repeat("a", 513):           "<invalid path>",
	} {
		if got := importDisplayPath(in); got != want {
			t.Errorf("importDisplayPath(%.40q) = %.40q, want %.40q", in, got, want)
		}
	}
}

func TestImporterRelativeKubeconfigEntry(t *testing.T) {
	setImporterEnv(t, "")

	cwd := t.TempDir()
	t.Chdir(cwd)
	for name, contents := range map[string]string{
		"kube/c.yaml":           plugintest.LoadFixture(t, "relative.yaml"),
		"kube/certs/client.crt": plugintest.LoadFixture(t, "client.crt"),
		"kube/certs/client.key": plugintest.LoadFixture(t, "client.key"),
		"~/.kube/config":        plugintest.LoadFixture(t, "literal-tilde.yaml"),
		"real/config":           plugintest.LoadFixture(t, "through-link.yaml"),
		"real/sub/.keep":        "",
		"config":                plugintest.LoadFixture(t, "lexical.yaml"),
	} {
		path := filepath.Join(cwd, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.Symlink(filepath.Join("real", "sub"), filepath.Join(cwd, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(os.DevNull, filepath.Join(cwd, "devnull")); err != nil {
		t.Fatal(err)
	}

	plugintest.TestImporter(t, TryKubeconfigFiles(), map[string]plugintest.ImportCase{
		"a literal ~/ KUBECONFIG entry is relative to the working directory": {
			Environment: map[string]string{
				"KUBECONFIG": "~/.kube/config",
			},
			Files: map[string]string{
				homeKubeconfig: plugintest.LoadFixture(t, "token.yaml"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile("~/.kube/config"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "literal")},
					},
				},
			},
		},
		"the entry is opened as written so symlinks resolve before ..": {
			Environment: map[string]string{
				"KUBECONFIG": "link/../config",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile("link/../config"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "through-link")},
					},
				},
			},
		},
		"a character device counts as an empty kubeconfig": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("devnull", "config"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile("config"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, "test-token-lexical", "lexical")},
					},
				},
			},
		},
		"different spellings that clean to the same path are different files": {
			Environment: map[string]string{
				"KUBECONFIG": kubeconfigList("link/../config", "config"),
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source:     importer.SourceFile("link/../config"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, importToken, "through-link")},
					},
					{
						Source:     importer.SourceFile("config"),
						Candidates: []sdk.ImportCandidate{tokenCandidate(prodImportAddr, "test-token-lexical", "lexical")},
					},
				},
			},
		},
		"I18 relative KUBECONFIG entry with references relative to that file": {
			Environment: map[string]string{
				"KUBECONFIG": "./kube/c.yaml",
			},
			ExpectedOutput: &sdk.ImportOutput{
				Attempts: []*sdk.ImportAttempt{
					{
						Source: importer.SourceFile("./kube/c.yaml"),
						Candidates: []sdk.ImportCandidate{
							{
								Fields: map[sdk.FieldName]string{
									fieldname.Address:     devImportAddr,
									fieldname.Certificate: base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, "client.crt"))),
									fieldname.PrivateKey:  base64.StdEncoding.EncodeToString([]byte(plugintest.LoadFixture(t, "client.key"))),
								},
								NameHint: "relative",
							},
						},
					},
				},
			},
		},
	})
}
