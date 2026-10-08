package argocd

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// addressAwareProvisioner skips provisioning when argocd targets a server other than the item's
// Address, since ARGOCD_SERVER and ARGOCD_AUTH_TOKEN override the selected argocd context.
type addressAwareProvisioner struct {
	sdk.Provisioner
}

func (p addressAwareProvisioner) Provision(ctx context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	if address := in.ItemFields[fieldname.Address]; address != "" {
		target := targetServer(in.HomeDir, out.CommandLine)
		if target != "" && normalizeServer(target) != normalizeServer(address) {
			return
		}
	}
	p.Provisioner.Provision(ctx, in, out)
}

type configFile struct {
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
}

func readConfig(homeDir string, args []string) *configFile {
	contents, err := os.ReadFile(configPath(homeDir, args))
	if err != nil {
		return nil
	}
	var config configFile
	if err := importer.FileContents(contents).ToYAML(&config); err != nil {
		return nil
	}
	return &config
}

// targetServer returns the server argocd would use without the plugin, or "" if unknown.
func targetServer(homeDir string, args []string) string {
	if server := flagValue(args, "--server"); server != "" {
		return server
	}
	if server := os.Getenv("ARGOCD_SERVER"); server != "" {
		return server
	}

	config := readConfig(homeDir, args)
	if config == nil {
		return ""
	}
	name := flagValue(args, "--argocd-context")
	if name == "" {
		name = config.CurrentContext
	}
	for _, c := range config.Contexts {
		if c.Name == name {
			return c.Server
		}
	}
	return ""
}

// configPath mirrors localconfig.DefaultConfigDir in argo-cd.
func configPath(homeDir string, args []string) string {
	if path := flagValue(args, "--config"); path != "" {
		return path
	}
	if dir := os.Getenv("ARGOCD_CONFIG_DIR"); dir != "" {
		return filepath.Join(dir, "config")
	}
	legacyDir := filepath.Join(homeDir, ".argocd")
	if _, err := os.Stat(legacyDir); err == nil {
		return filepath.Join(legacyDir, "config")
	}
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "argocd", "config")
	}
	return filepath.Join(homeDir, ".config", "argocd", "config")
}

// flagValue returns the value of the last "--flag value" or "--flag=value" before "--".
func flagValue(args []string, flag string) string {
	var value string
	for i := 0; i < len(args) && args[i] != "--"; i++ {
		if args[i] == flag && i+1 < len(args) {
			i++
			value = args[i]
		} else if v, ok := strings.CutPrefix(args[i], flag+"="); ok {
			value = v
		}
	}
	return value
}

func normalizeServer(server string) string {
	server = strings.ToLower(strings.TrimSpace(server))
	server = strings.TrimPrefix(server, "https://")
	server = strings.TrimPrefix(server, "http://")
	server = strings.TrimSuffix(server, "/")
	return strings.TrimSuffix(server, ":443")
}
