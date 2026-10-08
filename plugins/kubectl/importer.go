package kubectl

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

const maxDisplayedPathLength = 512

type kubeconfigSource struct {
	original string
	open     string
}

type importDiagnostic struct {
	key string
	err error
}

type materialSpec struct {
	kind        string
	inlineLabel string
	blockType   string
	data        string
	path        string
	entryName   string
	entrySource string
}

// TryKubeconfigFiles offers one candidate per context of the kubeconfig kubectl would load:
// the KUBECONFIG files when that variable is set, otherwise ~/.kube/config.
func TryKubeconfigFiles() sdk.Importer {
	return func(ctx context.Context, in sdk.ImportInput, out *sdk.ImportOutput) {
		attempts := map[string]*sdk.ImportAttempt{}
		spellings := map[string]string{}
		var files []*kubeconfig
		failed := false

		for _, src := range kubeconfigSources(in) {
			if info, err := os.Stat(src.open); errors.Is(err, os.ErrNotExist) || (err == nil && info.Mode()&os.ModeCharDevice != 0) {
				continue
			}
			k, err := readKubeconfig(src.open)
			if errors.Is(err, errNotFound) {
				continue
			}
			if err != nil {
				attempt := out.NewAttempt(importer.SourceFile(src.original))
				if errors.Is(err, errInvalidKubeconfig) {
					attempt.AddError(fmt.Errorf("parsing %s: not a valid kubeconfig", importDisplayPath(src.original)))
				} else {
					attempt.AddError(fmt.Errorf("reading %s: not a readable kubeconfig", importDisplayPath(src.original)))
				}
				failed = true
				continue
			}
			attempts[src.open] = out.NewAttempt(importer.SourceFile(src.original))
			spellings[src.open] = src.original
			files = append(files, k)
		}
		if failed || len(files) == 0 {
			return
		}

		merged := mergeKubeconfigs(files)
		var added []sdk.ImportCandidate
		reported := map[string]bool{}
		for _, kctx := range merged.Contexts {
			cl, ok := merged.cluster(kctx.Context.Cluster)
			if !ok || cl.Cluster.Server == "" {
				continue
			}
			u, ok := merged.user(kctx.Context.User)
			if !ok {
				continue
			}

			var diags []importDiagnostic
			fields, ok := importUserSecrets(in, kctx.Name, u, spellings[u.source], &diags)
			if ok {
				fields[fieldname.Address] = cl.Cluster.Server
				ca, caOK, diag := importMaterial(in, kctx.Name, spellings[cl.source], materialSpec{
					kind:        "certificate authority",
					inlineLabel: "certificate-authority-data",
					blockType:   "CERTIFICATE",
					data:        cl.Cluster.CertificateAuthorityData,
					path:        cl.Cluster.CertificateAuthority,
					entryName:   cl.Name,
					entrySource: cl.source,
				})
				if caOK {
					fields[fieldname.CertificateAuthority] = ca
				}
				if diag != nil {
					diags = append(diags, *diag)
				}
			}

			candidate := sdk.ImportCandidate{
				Fields:   fields,
				NameHint: importer.SanitizeNameHint(kctx.Name),
			}
			if ok && isDuplicateCandidate(added, candidate) {
				continue
			}

			attempt := attempts[kctx.source]
			for _, d := range diags {
				if !reported[d.key] {
					reported[d.key] = true
					attempt.AddError(d.err)
				}
			}
			if ok {
				added = append(added, candidate)
				attempt.AddCandidate(candidate)
			}
		}
	}
}

func kubeconfigSources(in sdk.ImportInput) []kubeconfigSource {
	var sources []kubeconfigSource
	if env := os.Getenv("KUBECONFIG"); env != "" {
		for _, entry := range filepath.SplitList(env) {
			if entry != "" {
				sources = append(sources, kubeconfigSource{original: entry, open: fromRootDir(in, entry)})
			}
		}
	} else {
		sources = append(sources, kubeconfigSource{original: "~/.kube/config", open: in.FromHomeDir(".kube", "config")})
	}

	seen := map[string]bool{}
	var unique []kubeconfigSource
	for _, src := range sources {
		if seen[src.open] {
			continue
		}
		seen[src.open] = true
		unique = append(unique, src)
	}
	return unique
}

func fromRootDir(in sdk.ImportInput, path string) string {
	if strings.HasPrefix(path, "/") {
		return strings.TrimSuffix(in.RootDir, "/") + path
	}
	return path
}

func resolveReference(in sdk.ImportInput, source, path string) string {
	if strings.HasPrefix(path, "/") {
		return fromRootDir(in, path)
	}
	return filepath.Join(filepath.Dir(source), path)
}

func importDisplayPath(path string) string {
	if strings.Contains(path, "\n") ||
		strings.HasPrefix(path, "-----BEGIN") ||
		strings.HasPrefix(path, "LS0t") ||
		len(path) > maxDisplayedPathLength {
		return "<invalid path>"
	}
	return path
}

func importUserSecrets(in sdk.ImportInput, contextName string, nu namedUser, spelling string, diags *[]importDiagnostic) (map[sdk.FieldName]string, bool) {
	u := nu.User
	if u.Exec != nil || u.AuthProvider != nil || u.TokenFile != "" {
		return nil, false
	}

	fields := map[sdk.FieldName]string{}
	if u.Token != "" {
		fields[fieldname.Token] = u.Token
	}

	cert, certOK, certDiag := importMaterial(in, contextName, spelling, materialSpec{
		kind:        "client certificate",
		inlineLabel: "client certificate",
		blockType:   "CERTIFICATE",
		data:        u.ClientCertificateData,
		path:        u.ClientCertificate,
		entryName:   nu.Name,
		entrySource: nu.source,
	})
	key, keyOK, keyDiag := importMaterial(in, contextName, spelling, materialSpec{
		kind:        "client key",
		inlineLabel: "client key",
		blockType:   "PRIVATE KEY",
		data:        u.ClientKeyData,
		path:        u.ClientKey,
		entryName:   nu.Name,
		entrySource: nu.source,
	})
	for _, d := range []*importDiagnostic{certDiag, keyDiag} {
		if d != nil {
			*diags = append(*diags, *d)
		}
	}
	if certOK && keyOK {
		fields[fieldname.Certificate] = cert
		fields[fieldname.PrivateKey] = key
	}

	if len(fields) == 0 {
		return nil, false
	}
	return fields, true
}

func importMaterial(in sdk.ImportInput, contextName, spelling string, m materialSpec) (string, bool, *importDiagnostic) {
	if m.data != "" {
		value, err := toBase64PEM(m.data, m.blockType)
		if err != nil {
			return "", false, &importDiagnostic{
				key: m.kind + "\x00inline\x00" + m.entrySource + "\x00" + m.entryName,
				err: fmt.Errorf("context %q: %s in %s could not be imported", contextName, m.inlineLabel, importDisplayPath(spelling)),
			}
		}
		return value, true, nil
	}
	if m.path == "" {
		return "", false, nil
	}

	resolved := resolveReference(in, m.entrySource, m.path)
	key := m.kind + "\x00file\x00" + filepath.Clean(resolved)
	contents, err := os.ReadFile(resolved)
	if err != nil {
		return "", false, &importDiagnostic{
			key: key,
			err: fmt.Errorf("context %q: %s at %s could not be imported", contextName, m.kind, importDisplayPath(m.path)),
		}
	}
	value, err := toBase64PEM(base64.StdEncoding.EncodeToString(contents), m.blockType)
	if err != nil {
		return "", false, &importDiagnostic{
			key: key,
			err: fmt.Errorf("context %q: %s in %s could not be imported", contextName, m.kind, importDisplayPath(m.path)),
		}
	}
	return value, true, nil
}

func isDuplicateCandidate(added []sdk.ImportCandidate, candidate sdk.ImportCandidate) bool {
	for _, existing := range added {
		if existing.Equal(candidate) {
			return true
		}
	}
	return false
}
