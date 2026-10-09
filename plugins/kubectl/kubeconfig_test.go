package kubectl

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const roundTripFixture = `apiVersion: v1
kind: Config
current-context: prod
preferences:
  colors: true
  extensions:
  - name: pref-ext
    extension:
      level: 3
extensions:
- name: top-ext
  extension:
    flag: true
clusters:
- name: prod
  cluster:
    server: https://prod.example.com:6443
    certificate-authority-data: Q0E=
    certificate-authority: /fake/ca.pem
    insecure-skip-tls-verify: true
    tls-server-name: api.internal.example.com
    proxy-url: http://proxy.example.com:3128
    disable-compression: true
    extensions:
    - name: cl-ext
      extension:
        region: eu
contexts:
- name: prod
  context:
    cluster: prod
    user: tokenuser
    namespace: team-a
    extensions:
    - name: ctx-ext
      extension:
        nested:
          list: [a, b]
users:
- name: tokenuser
  user:
    token: fake-token
    tokenFile: /fake/token
    as: someone
    as-uid: uid-1
    as-groups: [g1, g2]
    as-user-extra:
      scopes: [s1, s2]
    extensions:
    - name: usr-ext
      extension:
        key: value
- name: certuser
  user:
    client-certificate-data: Q0VSVA==
    client-key-data: S0VZ
    client-certificate: /fake/c.pem
    client-key: /fake/k.pem
    username: fake-user
    password: fake-pass
    exec:
      apiVersion: client.authentication.k8s.io/v1
      command: fake-cmd
      args: [one, two]
    auth-provider:
      name: oidc
      config:
        idp-issuer-url: https://idp.example.com
`

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func stripSources(k *kubeconfig) {
	for i := range k.Clusters {
		k.Clusters[i].source = ""
	}
	for i := range k.Contexts {
		k.Contexts[i].source = ""
	}
	for i := range k.Users {
		k.Users[i].source = ""
	}
}

func TestReadKubeconfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := writeFile(t, dir, "config", roundTripFixture)

	first, err := readKubeconfig(path)
	require.NoError(t, err)

	require.Equal(t, "prod", first.CurrentContext)
	require.Equal(t, true, first.Preferences["colors"])
	require.Len(t, first.Clusters, 1)
	require.True(t, first.Clusters[0].Cluster.InsecureSkipTLSVerify)
	require.Equal(t, "api.internal.example.com", first.Clusters[0].Cluster.TLSServerName)
	require.Equal(t, "http://proxy.example.com:3128", first.Clusters[0].Cluster.ProxyURL)
	require.True(t, first.Clusters[0].Cluster.DisableCompression)
	require.Len(t, first.Clusters[0].Cluster.Extensions, 1)
	require.Len(t, first.Extensions, 1)
	require.Equal(t, "team-a", first.Contexts[0].Context.Namespace)
	require.Len(t, first.Contexts[0].Context.Extensions, 1)
	require.Len(t, first.Users, 2)
	tokenUser := first.Users[0].User
	require.Equal(t, "fake-token", tokenUser.Token)
	require.Equal(t, "someone", tokenUser.As)
	require.Equal(t, "uid-1", tokenUser.AsUID)
	require.Equal(t, []string{"g1", "g2"}, tokenUser.AsGroups)
	require.Equal(t, map[string][]string{"scopes": {"s1", "s2"}}, tokenUser.AsUserExtra)
	require.Len(t, tokenUser.Extensions, 1)
	certUser := first.Users[1].User
	require.Equal(t, "Q0VSVA==", certUser.ClientCertificateData)
	require.Equal(t, "S0VZ", certUser.ClientKeyData)
	require.Equal(t, "fake-cmd", certUser.Exec["command"])
	require.Equal(t, "oidc", certUser.AuthProvider["name"])

	out, err := yaml.Marshal(first)
	require.NoError(t, err)
	for _, want := range []string{
		"as-uid: uid-1", "as-groups:", "as-user-extra:", "extensions:", "preferences:",
		"namespace: team-a", "insecure-skip-tls-verify: true", "tokenFile: /fake/token",
		"client-certificate-data:", "auth-provider:", "exec:", "idp-issuer-url",
		"tls-server-name:", "proxy-url:", "disable-compression: true", "top-ext", "cl-ext",
	} {
		require.Contains(t, string(out), want)
	}

	path2 := writeFile(t, dir, "config2", string(out))
	second, err := readKubeconfig(path2)
	require.NoError(t, err)

	stripSources(first)
	stripSources(second)
	require.Equal(t, first, second)
}

func TestReadKubeconfigSourceAndErrors(t *testing.T) {
	dir := t.TempDir()

	t.Run("source is set on every entry", func(t *testing.T) {
		path := writeFile(t, dir, "config", roundTripFixture)
		k, err := readKubeconfig(path)
		require.NoError(t, err)
		require.Equal(t, path, k.Clusters[0].source)
		require.Equal(t, path, k.Contexts[0].source)
		require.Equal(t, path, k.Users[0].source)
		require.Equal(t, path, k.Users[1].source)
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := readKubeconfig(filepath.Join(dir, "nope"))
		require.True(t, errors.Is(err, errNotFound))
	})

	t.Run("directory", func(t *testing.T) {
		_, err := readKubeconfig(dir)
		require.Error(t, err)
		require.False(t, errors.Is(err, errNotFound))
	})

	t.Run("fifo", func(t *testing.T) {
		fifo := filepath.Join(dir, "fifo")
		if err := exec.Command("mkfifo", fifo).Run(); err != nil {
			t.Skip("mkfifo unavailable")
		}
		_, err := readKubeconfig(fifo)
		require.Error(t, err)
		require.False(t, errors.Is(err, errNotFound))
	})

	t.Run("duplicate names", func(t *testing.T) {
		for kind, content := range map[string]string{
			"clusters": "clusters:\n- name: a\n  cluster: {server: 'https://x'}\n- name: a\n  cluster: {server: 'https://y'}\n",
			"contexts": "contexts:\n- name: a\n  context: {cluster: c, user: u}\n- name: a\n  context: {cluster: c, user: u}\n",
			"users":    "users:\n- name: a\n  user: {token: fake-secret-token}\n- name: a\n  user: {token: other}\n",
		} {
			path := writeFile(t, dir, "dup-"+kind, content)
			_, err := readKubeconfig(path)
			require.Error(t, err, kind)
			require.NotContains(t, err.Error(), "fake-secret-token")
		}
	})

	t.Run("invalid yaml never leaks content", func(t *testing.T) {
		path := writeFile(t, dir, "bad", "clusters: [fake-secret-token: {\n")
		_, err := readKubeconfig(path)
		require.Error(t, err)
		require.Equal(t, "parsing "+path+": not a valid kubeconfig", err.Error())

		path = writeFile(t, dir, "bad-type", "clusters: fake-secret-token\n")
		_, err = readKubeconfig(path)
		require.Equal(t, "parsing "+path+": not a valid kubeconfig", err.Error())
	})
}

func TestLoadKubeconfigMerge(t *testing.T) {
	dir := t.TempDir()
	a := writeFile(t, dir, "a", `current-context: ""
contexts:
- name: prod
  context: {cluster: prod, user: alice}
clusters:
- name: other
  cluster: {server: 'https://127.0.0.1:7443'}
preferences:
  colors: true
`)
	b := writeFile(t, dir, "b", `current-context: prod
contexts:
- name: prod
  context: {cluster: prod, user: bob}
clusters:
- name: prod
  cluster: {server: 'https://prod.example.com:6443'}
- name: other
  cluster: {server: 'https://shadowed.example.com'}
users:
- name: bob
  user: {token: fake}
preferences:
  colors: false
`)
	missing := filepath.Join(dir, "missing")

	k, err := loadKubeconfig([]string{missing, a, missing, b})
	require.NoError(t, err)
	require.Equal(t, "prod", k.CurrentContext)

	ctx, ok := k.context("prod")
	require.True(t, ok)
	require.Equal(t, "alice", ctx.Context.User)
	require.Equal(t, a, ctx.source)

	cl, ok := k.cluster("prod")
	require.True(t, ok)
	require.Equal(t, "https://prod.example.com:6443", cl.Cluster.Server)
	require.Equal(t, b, cl.source)

	other, ok := k.cluster("other")
	require.True(t, ok)
	require.Equal(t, "https://127.0.0.1:7443", other.Cluster.Server)
	require.Equal(t, a, other.source)

	u, ok := k.user("bob")
	require.True(t, ok)
	require.Equal(t, b, u.source)

	_, ok = k.context("nope")
	require.False(t, ok)
	_, ok = k.cluster("nope")
	require.False(t, ok)
	_, ok = k.user("nope")
	require.False(t, ok)

	require.Equal(t, true, k.Preferences["colors"])

	t.Run("all missing", func(t *testing.T) {
		k, err := loadKubeconfig([]string{missing, filepath.Join(dir, "missing2")})
		require.Nil(t, k)
		require.True(t, errors.Is(err, errNotFound))

		k, err = loadKubeconfig(nil)
		require.Nil(t, k)
		require.True(t, errors.Is(err, errNotFound))
	})

	t.Run("parse failure in any present file", func(t *testing.T) {
		dup := writeFile(t, dir, "dup", "contexts:\n- name: x\n- name: x\n")
		_, err := loadKubeconfig([]string{a, dup, b})
		require.Error(t, err)
		require.False(t, errors.Is(err, errNotFound))
	})

	t.Run("directory in the list", func(t *testing.T) {
		_, err := loadKubeconfig([]string{a, dir})
		require.Error(t, err)
		require.False(t, errors.Is(err, errNotFound))
	})
}

func TestMergeKubeconfigs(t *testing.T) {
	first := &kubeconfig{CurrentContext: ""}
	second := &kubeconfig{CurrentContext: "x", Preferences: map[string]any{"colors": true}, Extensions: []any{"e2"}}
	third := &kubeconfig{CurrentContext: "y", Extensions: []any{"e3"}}
	merged := mergeKubeconfigs([]*kubeconfig{first, nil, second, third})
	require.Equal(t, "x", merged.CurrentContext)
	require.Equal(t, map[string]any{"colors": true}, merged.Preferences)
	require.Equal(t, []any{"e2"}, merged.Extensions)
}

func TestHasName(t *testing.T) {
	k := &kubeconfig{
		CurrentContext: "current",
		Clusters:       []namedCluster{{Name: "c1"}},
		Contexts:       []namedContext{{Name: "ctx1", Context: kubeContext{Cluster: "refc", User: "refu"}}},
		Users:          []namedUser{{Name: "u1"}},
	}
	for _, name := range []string{"c1", "ctx1", "u1", "refc", "refu", "current"} {
		require.True(t, k.hasName(name), name)
	}
	require.False(t, k.hasName("absent"))
	require.False(t, k.hasName(""))
}

func TestNormalizeServer(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"HTTPS://Prod.Example.com:6443/", "https://prod.example.com:6443", true},
		{"  https://prod.example.com:6443  ", "https://prod.example.com:6443", true},
		{"https://x:443", "https://x", true},
		{"http://x:80", "http://x", true},
		{"http://x:443", "http://x:443", true},
		{"https://x:80", "https://x:80", true},
		{"prod.example.com:6443", "https://prod.example.com:6443", true},
		{"prod.example.com", "https://prod.example.com", true},
		{"127.0.0.1:6443", "https://127.0.0.1:6443", true},
		{"https://rancher.example.com/k8s/clusters/c-abc/", "https://rancher.example.com/k8s/clusters/c-abc", true},
		{"https://x/a?b=c#d", "https://x/a", true},
		{"https://x:6443?b=c", "https://x:6443", true},
		{"https://[fe80::1%25en0]:6443", "https://[fe80::1%25en0]:6443", true},
		{"https://[fe80::1%25ETH0]:6443", "https://[fe80::1%25ETH0]:6443", true},
		{"example.com/path?redirect=https://other", "https://example.com/path", true},
		{"example.com/a://b", "https://example.com/a://b", true},
		{"example.com#frag://x", "https://example.com", true},
		{"https://[fe80::1%25en%25x]", "https://[fe80::1%25en%25x]", true},
		{"https://[fe80::1%25a%20b]:6443", "https://[fe80::1%25a%20b]:6443", true},
		{"https://[fe80::1%25a%252Fb]", "https://[fe80::1%25a%252Fb]", true},
		{"https://[fe80::1%25a%2Fb]:6443", "", false},
		{"https://[fe80::1%25]", "", false},
		{"https://[::1]:443", "https://[::1]", true},
		{"https://[::1]:6443", "https://[::1]:6443", true},
		{"http://[::1]:80/p/", "http://[::1]/p", true},
		{"https://user:pw@x", "", false},
		{"https://user@x", "", false},
		{"ftp://x", "", false},
		{"https://", "", false},
		{"https://:443", "", false},
		{"::::", "", false},
		{"", "", false},
		{"   ", "", false},
		{"https://a b", "", false},
		{"https://x:notaport", "", false},
		{"https://[notanip]:6443", "", false},
		{"https://[1.2.3.4]", "", false},
	}
	for _, tt := range tests {
		got, ok := normalizeServer(tt.in)
		require.Equal(t, tt.ok, ok, tt.in)
		require.Equal(t, tt.want, got, tt.in)
		if ok {
			_, err := url.Parse(got)
			require.NoError(t, err, got)
			again, ok := normalizeServer(got)
			require.True(t, ok, got)
			require.Equal(t, got, again, got)
		}
	}

	httpSrv, _ := normalizeServer("http://x")
	httpsSrv, _ := normalizeServer("https://x")
	require.NotEqual(t, httpSrv, httpsSrv)
}

func throwawayPEMs(t *testing.T) (certPEM, pkcs8PEM, ecPEM, rsaPEM string) {
	t.Helper()
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "throwaway"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &ecKey.PublicKey, ecKey)
	require.NoError(t, err)
	certPEM = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))

	pkcs8, err := x509.MarshalPKCS8PrivateKey(ecKey)
	require.NoError(t, err)
	pkcs8PEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8}))

	ecDER, err := x509.MarshalECPrivateKey(ecKey)
	require.NoError(t, err)
	ecPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: ecDER}))

	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	require.NoError(t, err)
	rsaPEM = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(rsaKey)}))
	return
}

func wrapBase64(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i += 20 {
		end := min(i+20, len(s))
		b.WriteString(s[i:end])
		if i%40 == 0 {
			b.WriteString("\n")
		} else {
			b.WriteString(" ")
		}
	}
	return "  " + b.String() + "\n"
}

func TestToBase64PEM(t *testing.T) {
	certPEM, pkcs8PEM, ecPEM, rsaPEM := throwawayPEMs(t)

	t.Run("pem and base64 of the same pem agree", func(t *testing.T) {
		want := base64.StdEncoding.EncodeToString([]byte(certPEM))

		got, err := toBase64PEM(certPEM, "CERTIFICATE")
		require.NoError(t, err)
		require.Equal(t, want, got)

		got, err = toBase64PEM(want, "CERTIFICATE")
		require.NoError(t, err)
		require.Equal(t, want, got)

		got, err = toBase64PEM(wrapBase64(want), "CERTIFICATE")
		require.NoError(t, err)
		require.Equal(t, want, got)

		for _, padded := range []string{"\n" + certPEM, certPEM + "\n\n  ", "\n\n" + certPEM + "  \n"} {
			got, err = toBase64PEM(padded, "CERTIFICATE")
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(padded)), got)

			viaBase64, err := toBase64PEM(base64.StdEncoding.EncodeToString([]byte(padded)), "CERTIFICATE")
			require.NoError(t, err)
			require.Equal(t, got, viaBase64)
		}
	})

	t.Run("private key block types", func(t *testing.T) {
		for _, p := range []string{pkcs8PEM, ecPEM, rsaPEM} {
			got, err := toBase64PEM(p, "PRIVATE KEY")
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(p)), got)

			got, err = toBase64PEM(base64.StdEncoding.EncodeToString([]byte(p)), "PRIVATE KEY")
			require.NoError(t, err)
			require.Equal(t, base64.StdEncoding.EncodeToString([]byte(p)), got)
		}
	})

	t.Run("block type mismatch", func(t *testing.T) {
		_, err := toBase64PEM(certPEM, "PRIVATE KEY")
		require.Error(t, err)
		_, err = toBase64PEM(pkcs8PEM, "CERTIFICATE")
		require.Error(t, err)
		_, err = toBase64PEM(base64.StdEncoding.EncodeToString([]byte(certPEM)), "PRIVATE KEY")
		require.Error(t, err)
	})

	t.Run("one of several types", func(t *testing.T) {
		_, err := toBase64PEM(certPEM, "PRIVATE KEY", "CERTIFICATE")
		require.NoError(t, err)
	})

	t.Run("invalid input never echoed", func(t *testing.T) {
		secretish := base64.StdEncoding.EncodeToString([]byte("fake-secret-token-not-pem"))
		for _, in := range []string{
			"not base64!",
			secretish,
			"-----BEGIN CERTIFICATE-----\nfake-secret-token\n-----END CERTIFICATE-----",
			"-----BEGIN",
			"",
		} {
			_, err := toBase64PEM(in, "CERTIFICATE")
			require.Error(t, err, in)
			if in != "" {
				require.NotContains(t, err.Error(), in)
			}
			require.NotContains(t, err.Error(), "fake-secret-token")
		}
		_, err := toBase64PEM("not base64!", "CERTIFICATE")
		require.EqualError(t, err, "value is neither PEM nor base64-encoded PEM")
	})
}

func TestCertificateAuthorityFieldName(t *testing.T) {
	require.Equal(t, "Certificate Authority", string(fieldname.CertificateAuthority))
	require.Contains(t, fieldname.ListAll(), fieldname.CertificateAuthority)
}
