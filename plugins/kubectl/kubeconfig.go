package kubectl

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/1Password/shell-plugins/sdk/importer"
)

type kubeconfig struct {
	APIVersion     string         `yaml:"apiVersion,omitempty"`
	Kind           string         `yaml:"kind,omitempty"`
	Preferences    map[string]any `yaml:"preferences,omitempty"`
	Clusters       []namedCluster `yaml:"clusters,omitempty"`
	Contexts       []namedContext `yaml:"contexts,omitempty"`
	Users          []namedUser    `yaml:"users,omitempty"`
	CurrentContext string         `yaml:"current-context,omitempty"`
	Extensions     []any          `yaml:"extensions,omitempty"`
}

type namedCluster struct {
	Name    string  `yaml:"name"`
	Cluster cluster `yaml:"cluster"`
	source  string
}

type cluster struct {
	Server                   string `yaml:"server,omitempty"`
	CertificateAuthority     string `yaml:"certificate-authority,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
	TLSServerName            string `yaml:"tls-server-name,omitempty"`
	ProxyURL                 string `yaml:"proxy-url,omitempty"`
	DisableCompression       bool   `yaml:"disable-compression,omitempty"`
	Extensions               []any  `yaml:"extensions,omitempty"`
}

type namedContext struct {
	Name    string      `yaml:"name"`
	Context kubeContext `yaml:"context"`
	source  string
}

type kubeContext struct {
	Cluster    string `yaml:"cluster"`
	User       string `yaml:"user"`
	Namespace  string `yaml:"namespace,omitempty"`
	Extensions []any  `yaml:"extensions,omitempty"`
}

type namedUser struct {
	Name   string `yaml:"name"`
	User   user   `yaml:"user"`
	source string
}

type user struct {
	Token                 string              `yaml:"token,omitempty"`
	TokenFile             string              `yaml:"tokenFile,omitempty"`
	ClientCertificate     string              `yaml:"client-certificate,omitempty"`
	ClientCertificateData string              `yaml:"client-certificate-data,omitempty"`
	ClientKey             string              `yaml:"client-key,omitempty"`
	ClientKeyData         string              `yaml:"client-key-data,omitempty"`
	Username              string              `yaml:"username,omitempty"`
	Password              string              `yaml:"password,omitempty"`
	Exec                  map[string]any      `yaml:"exec,omitempty"`
	AuthProvider          map[string]any      `yaml:"auth-provider,omitempty"`
	As                    string              `yaml:"as,omitempty"`
	AsUID                 string              `yaml:"as-uid,omitempty"`
	AsGroups              []string            `yaml:"as-groups,omitempty"`
	AsUserExtra           map[string][]string `yaml:"as-user-extra,omitempty"`
	Extensions            []any               `yaml:"extensions,omitempty"`
}

var errNotFound = errors.New("kubeconfig file not found")

var errInvalidKubeconfig = errors.New("not a valid kubeconfig")

func invalidKubeconfig(path string) error {
	return fmt.Errorf("parsing %s: %w", path, errInvalidKubeconfig)
}

func readKubeconfig(path string) (*kubeconfig, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("reading %s: not a regular file", path)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var k kubeconfig
	if err := importer.FileContents(contents).ToYAML(&k); err != nil {
		return nil, invalidKubeconfig(path)
	}

	clusters := map[string]bool{}
	for i := range k.Clusters {
		if clusters[k.Clusters[i].Name] {
			return nil, invalidKubeconfig(path)
		}
		clusters[k.Clusters[i].Name] = true
		k.Clusters[i].source = path
	}
	contexts := map[string]bool{}
	for i := range k.Contexts {
		if contexts[k.Contexts[i].Name] {
			return nil, invalidKubeconfig(path)
		}
		contexts[k.Contexts[i].Name] = true
		k.Contexts[i].source = path
	}
	users := map[string]bool{}
	for i := range k.Users {
		if users[k.Users[i].Name] {
			return nil, invalidKubeconfig(path)
		}
		users[k.Users[i].Name] = true
		k.Users[i].source = path
	}

	return &k, nil
}

func mergeKubeconfigs(files []*kubeconfig) *kubeconfig {
	merged := &kubeconfig{}
	clusters := map[string]bool{}
	contexts := map[string]bool{}
	users := map[string]bool{}

	for _, f := range files {
		if f == nil {
			continue
		}
		if merged.APIVersion == "" {
			merged.APIVersion = f.APIVersion
		}
		if merged.Kind == "" {
			merged.Kind = f.Kind
		}
		if merged.CurrentContext == "" {
			merged.CurrentContext = f.CurrentContext
		}
		if len(merged.Preferences) == 0 {
			merged.Preferences = f.Preferences
		}
		if len(merged.Extensions) == 0 {
			merged.Extensions = f.Extensions
		}
		for _, c := range f.Clusters {
			if !clusters[c.Name] {
				clusters[c.Name] = true
				merged.Clusters = append(merged.Clusters, c)
			}
		}
		for _, c := range f.Contexts {
			if !contexts[c.Name] {
				contexts[c.Name] = true
				merged.Contexts = append(merged.Contexts, c)
			}
		}
		for _, u := range f.Users {
			if !users[u.Name] {
				users[u.Name] = true
				merged.Users = append(merged.Users, u)
			}
		}
	}

	return merged
}

func loadKubeconfig(paths []string) (*kubeconfig, error) {
	var files []*kubeconfig
	for _, path := range paths {
		k, err := readKubeconfig(path)
		if errors.Is(err, errNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		files = append(files, k)
	}
	if len(files) == 0 {
		return nil, errNotFound
	}

	return mergeKubeconfigs(files), nil
}

func (k *kubeconfig) context(name string) (namedContext, bool) {
	for _, c := range k.Contexts {
		if c.Name == name {
			return c, true
		}
	}
	return namedContext{}, false
}

func (k *kubeconfig) cluster(name string) (namedCluster, bool) {
	for _, c := range k.Clusters {
		if c.Name == name {
			return c, true
		}
	}
	return namedCluster{}, false
}

func (k *kubeconfig) user(name string) (namedUser, bool) {
	for _, u := range k.Users {
		if u.Name == name {
			return u, true
		}
	}
	return namedUser{}, false
}

func (k *kubeconfig) hasName(name string) bool {
	if k.CurrentContext == name {
		return true
	}
	for _, c := range k.Clusters {
		if c.Name == name {
			return true
		}
	}
	for _, u := range k.Users {
		if u.Name == name {
			return true
		}
	}
	for _, c := range k.Contexts {
		if c.Name == name || c.Context.Cluster == name || c.Context.User == name {
			return true
		}
	}
	return false
}

var schemePrefix = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+.-]*://`)

func normalizeServer(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if !schemePrefix.MatchString(s) {
		s = "https://" + s
	}

	u, err := url.Parse(s)
	if err != nil || u.User != nil {
		return "", false
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", false
	}

	host := u.Hostname()
	if host == "" {
		return "", false
	}
	if strings.HasPrefix(u.Host, "[") {
		addr, err := netip.ParseAddr(host)
		if err != nil || !addr.Is6() {
			return "", false
		}
		address, zone, hasZone := strings.Cut(host, "%")
		host = address
		if hasZone {
			host += "%25" + escapeZone(zone)
		}
		host = "[" + host + "]"
	} else {
		if strings.Contains(host, ":") {
			return "", false
		}
		host = strings.ToLower(host)
	}

	port := u.Port()
	if (scheme == "https" && port == "443") || (scheme == "http" && port == "80") {
		port = ""
	}
	if port != "" {
		host += ":" + port
	}

	normalized := scheme + "://" + host + strings.TrimRight(u.EscapedPath(), "/")
	if _, err := url.Parse(normalized); err != nil {
		return "", false
	}

	return normalized, true
}

func escapeZone(zone string) string {
	var b strings.Builder
	for i := 0; i < len(zone); i++ {
		c := zone[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '.', c == '_', c == '~':
			b.WriteByte(c)
		default:
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func toBase64PEM(value string, blockTypes ...string) (string, error) {
	var pemBytes []byte
	if strings.HasPrefix(strings.TrimLeftFunc(value, unicode.IsSpace), "-----BEGIN") {
		pemBytes = []byte(value)
	} else {
		compact := strings.Join(strings.FieldsFunc(value, unicode.IsSpace), "")
		decoded, err := base64.StdEncoding.DecodeString(compact)
		if err != nil {
			return "", errors.New("value is neither PEM nor base64-encoded PEM")
		}
		pemBytes = decoded
	}

	block, _ := pem.Decode(bytes.TrimLeftFunc(pemBytes, unicode.IsSpace))
	if block == nil {
		return "", errors.New("value does not contain a PEM block")
	}
	if len(blockTypes) > 0 {
		matched := false
		for _, t := range blockTypes {
			if strings.HasSuffix(block.Type, t) {
				matched = true
				break
			}
		}
		if !matched {
			return "", fmt.Errorf("PEM block is not of type %s", strings.Join(blockTypes, " or "))
		}
	}

	return base64.StdEncoding.EncodeToString(pemBytes), nil
}
