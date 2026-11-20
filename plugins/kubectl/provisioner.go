package kubectl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
	"gopkg.in/yaml.v2"
)

const (
	overlayNamePrefix = "1password-shell-plugin-"
	maxNameAttempts   = 16
)

var (
	errNameGeneration = errors.New("kubectl: could not generate a name for the kubeconfig overlay")
	errNameTaken      = errors.New("kubectl: could not choose a unique name for the kubeconfig overlay")
	randRead          = rand.Read
)

var explicitAuthFlags = []string{"token", "client-certificate", "client-key", "username", "password", "user"}

type kubeconfigProvisioner struct {
	newName func() string
}

func Provisioner() sdk.Provisioner {
	return kubeconfigProvisioner{}
}

func (p kubeconfigProvisioner) name() (string, error) {
	if p.newName != nil {
		return p.newName(), nil
	}
	b := make([]byte, 8)
	if _, err := randRead(b); err != nil {
		return "", errNameGeneration
	}
	return overlayNamePrefix + hex.EncodeToString(b), nil
}

func (p kubeconfigProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if len(out.CommandLine) == 0 {
		return
	}
	args := parseArgs(out.CommandLine[1:])
	sep := string(os.PathListSeparator)

	if args.isAmbiguous() {
		return
	}
	for _, flag := range explicitAuthFlags {
		if args.has(flag) {
			return
		}
	}

	envKubeconfig := os.Getenv("KUBECONFIG")
	homeConfig := in.FromHomeDir(".kube", "config")
	explicitPath, explicit := args.value("kubeconfig")
	var sources []string
	switch {
	case explicit:
		if strings.Contains(explicitPath, sep) {
			return
		}
		sources = []string{explicitPath}
	case envKubeconfig != "":
		for _, path := range filepath.SplitList(envKubeconfig) {
			if path != "" {
				sources = append(sources, path)
			}
		}
	default:
		if strings.Contains(homeConfig, sep) {
			return
		}
		sources = []string{homeConfig}
	}

	if kubercInterferes(in.HomeDir, args) {
		return
	}

	if !explicit && envKubeconfig == "" && isFreshMachine(in, args) {
		p.provisionStandalone(in, args, out)
		return
	}

	merged, err := loadKubeconfig(sources)
	if err != nil {
		return
	}

	ctxName, ok := args.value("context")
	if !ok {
		ctxName = merged.CurrentContext
	}
	target, ok := merged.context(ctxName)
	if ctxName == "" || !ok {
		return
	}
	clusterName, ok := args.value("cluster")
	if !ok {
		clusterName = target.Context.Cluster
	}
	server, ok := args.value("server")
	if !ok {
		c, _ := merged.cluster(clusterName)
		server = c.Cluster.Server
	}
	if server == "" {
		return
	}

	address := in.ItemFields[fieldname.Address]
	if address == "" {
		out.AddError(errors.New("kubectl: the 1Password item has no Address; set it to the cluster server URL"))
		return
	}
	if !sameServer(server, address) {
		return
	}

	credentials, errs := itemUser(in.ItemFields)
	if len(errs) > 0 {
		for _, err := range errs {
			out.AddError(err)
		}
		return
	}
	if original, ok := merged.user(target.Context.User); ok {
		credentials.As = original.User.As
		credentials.AsUID = original.User.AsUID
		credentials.AsGroups = original.User.AsGroups
		credentials.AsUserExtra = original.User.AsUserExtra
		credentials.Extensions = original.User.Extensions
	}

	name, err := p.uniqueName(merged)
	if err != nil {
		out.AddError(err)
		return
	}

	overlay := kubeconfig{
		APIVersion: "v1",
		Kind:       "Config",
		Contexts: []namedContext{{
			Name: ctxName,
			Context: kubeContext{
				Cluster:    target.Context.Cluster,
				User:       name,
				Namespace:  target.Context.Namespace,
				Extensions: target.Context.Extensions,
			},
		}},
		Users: []namedUser{{Name: name, User: credentials}},
	}

	var rest string
	switch {
	case explicit:
		rest = explicitPath
	case envKubeconfig != "":
		rest = envKubeconfig
	default:
		rest = homeConfig
	}
	if !writeOverlay(in, out, overlay, rest) {
		return
	}
	if explicit {
		out.CommandLine = append([]string{out.CommandLine[0]}, removeFlag(out.CommandLine[1:], "kubeconfig")...)
	}
}

func (p kubeconfigProvisioner) Deprovision(ctx context.Context, in sdk.DeprovisionInput, out *sdk.DeprovisionOutput) {
	// The SDK removes the temp dir holding the overlay.
}

func (p kubeconfigProvisioner) Description() string {
	return "Provision a temporary kubeconfig overlay for the context that targets the item's cluster"
}

func (p kubeconfigProvisioner) uniqueName(merged *kubeconfig) (string, error) {
	for i := 0; i < maxNameAttempts; i++ {
		name, err := p.name()
		if err != nil {
			return "", err
		}
		if !merged.hasName(name) {
			return name, nil
		}
	}
	return "", errNameTaken
}

// isFreshMachine reports whether kubectl has no kubeconfig and no in-cluster config to fall
// back on, so that it would otherwise contact localhost:8080.
func isFreshMachine(in sdk.ProvisionInput, args parsedArgs) bool {
	for _, path := range []string{in.FromHomeDir(".kube", "config"), in.FromHomeDir(".kube", ".kubeconfig")} {
		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" || os.Getenv("KUBERNETES_MASTER") != "" {
		return false
	}
	return !args.has("context") && !args.has("cluster") && !args.has("kubeconfig")
}

func (p kubeconfigProvisioner) provisionStandalone(in sdk.ProvisionInput, args parsedArgs, out *sdk.ProvisionOutput) {
	address := strings.TrimSpace(in.ItemFields[fieldname.Address])
	if address == "" {
		out.AddError(errors.New("kubectl: the 1Password item has no Address; set it to the cluster server URL"))
		return
	}
	if _, ok := normalizeServer(address); !ok {
		out.AddError(errors.New("kubectl: the 1Password item's Address is not a valid https:// or http:// server URL"))
		return
	}
	c := cluster{Server: address}
	if !schemePrefix.MatchString(address) {
		c.Server = "https://" + address
	}
	if server, ok := args.value("server"); ok && !sameServer(server, c.Server) {
		return
	}

	credentials, errs := itemUser(in.ItemFields)
	if ca := in.ItemFields[fieldname.CertificateAuthority]; ca != "" {
		data, err := toBase64PEM(ca, "CERTIFICATE")
		if err != nil {
			errs = append(errs, fieldError(fieldname.CertificateAuthority, err))
		}
		c.CertificateAuthorityData = data
	}
	if len(errs) > 0 {
		for _, err := range errs {
			out.AddError(err)
		}
		return
	}

	name, err := p.name()
	if err != nil {
		out.AddError(err)
		return
	}
	standalone := kubeconfig{
		APIVersion:     "v1",
		Kind:           "Config",
		Clusters:       []namedCluster{{Name: name, Cluster: c}},
		Contexts:       []namedContext{{Name: name, Context: kubeContext{Cluster: name, User: name}}},
		Users:          []namedUser{{Name: name, User: credentials}},
		CurrentContext: name,
	}
	writeOverlay(in, out, standalone, "")
}

// writeOverlay adds the overlay file and points KUBECONFIG at it, followed by rest when set.
func writeOverlay(in sdk.ProvisionInput, out *sdk.ProvisionOutput, config kubeconfig, rest string) bool {
	sep := string(os.PathListSeparator)
	path := in.FromTempDir("kubeconfig")
	if strings.Contains(path, sep) {
		return false
	}
	contents, err := yaml.Marshal(config)
	if err != nil {
		out.AddError(errors.New("kubectl: could not build the kubeconfig overlay"))
		return false
	}

	out.AddSecretFile(path, contents)
	if rest != "" {
		path += sep + rest
	}
	out.AddEnvVar("KUBECONFIG", path)
	return true
}

func itemUser(fields map[sdk.FieldName]string) (user, []error) {
	var u user
	var errs []error
	u.Token = strings.TrimSpace(fields[fieldname.Token])

	cert := fields[fieldname.Certificate]
	key := fields[fieldname.PrivateKey]
	switch {
	case (cert == "") != (key == ""):
		errs = append(errs, errors.New("kubectl: the 1Password item needs both Certificate and Private Key"))
	case cert != "":
		var err error
		if u.ClientCertificateData, err = toBase64PEM(cert, "CERTIFICATE"); err != nil {
			errs = append(errs, fieldError(fieldname.Certificate, err))
		}
		if u.ClientKeyData, err = toBase64PEM(key, "PRIVATE KEY"); err != nil {
			errs = append(errs, fieldError(fieldname.PrivateKey, err))
		}
	case u.Token == "":
		errs = append(errs, errors.New("kubectl: the 1Password item has neither Token nor Certificate and Private Key"))
	}
	return u, errs
}

// fieldError wraps a toBase64PEM error, whose message never contains the value.
func fieldError(field sdk.FieldName, err error) error {
	return fmt.Errorf("kubectl: %s in the 1Password item: %w", field, err)
}

// sameServer requires both sides to be scheme-less or both to carry a scheme: kubectl
// talks plain http to a scheme-less server that has no TLS material, so "host:port" must
// not match "https://host:port".
func sameServer(a, b string) bool {
	a, b = strings.TrimSpace(a), strings.TrimSpace(b)
	if schemePrefix.MatchString(a) != schemePrefix.MatchString(b) {
		return false
	}
	na, okA := normalizeServer(a)
	nb, okB := normalizeServer(b)
	if !okA || !okB || na != nb {
		return false
	}
	// Without a scheme kubectl may use http, so the :443 default that normalizeServer
	// drops is not a default here; the written host:port must agree as well.
	if !schemePrefix.MatchString(a) {
		return strings.EqualFold(strings.TrimSuffix(a, "/"), strings.TrimSuffix(b, "/"))
	}
	return true
}
