package kubectl

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

type kubercOption struct {
	Name string `yaml:"name"`
}

type kubercCommandDefaults struct {
	Options []kubercOption `yaml:"options"`
	Flags   []kubercOption `yaml:"flags"`
}

type kubercAlias struct {
	Name string `yaml:"name"`
}

// kubercPreference covers v1beta1 (defaults/options) and v1alpha1 (overrides/flags).
type kubercPreference struct {
	APIVersion string                  `yaml:"apiVersion"`
	Kind       string                  `yaml:"kind"`
	Aliases    []kubercAlias           `yaml:"aliases"`
	Defaults   []kubercCommandDefaults `yaml:"defaults"`
	Overrides  []kubercCommandDefaults `yaml:"overrides"`
}

var kubercTargetingOptions = map[string]bool{
	"context":                  true,
	"cluster":                  true,
	"server":                   true,
	"kubeconfig":               true,
	"user":                     true,
	"token":                    true,
	"client-certificate":       true,
	"client-key":               true,
	"username":                 true,
	"password":                 true,
	"proxy-url":                true,
	"tls-server-name":          true,
	"insecure-skip-tls-verify": true,
	"certificate-authority":    true,
}

// kubercInterferes reports whether a kuberc file in effect could change which cluster or
// credentials the command uses. kubectl applies defaults through cobra command aliases and
// parseArgs cannot see subcommand-local flag values, so the check ignores which command an
// entry is for and fails closed.
func kubercInterferes(homeDir string, args parsedArgs) bool {
	if len(args.flags["kuberc"]) > 1 {
		return true
	}
	path, ok := kubercPath(homeDir, args)
	if !ok {
		return false
	}

	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err != nil || !info.Mode().IsRegular() {
		return true
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		return true
	}

	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	for {
		var pref kubercPreference
		err := decoder.Decode(&pref)
		if errors.Is(err, io.EOF) {
			return false
		}
		if err != nil || pref.interferes(args) {
			return true
		}
	}
}

func kubercPath(homeDir string, args parsedArgs) (string, bool) {
	if path, ok := args.value("kuberc"); ok {
		return path, true
	}
	if env := os.Getenv("KUBERC"); env != "" {
		return env, env != "off"
	}
	return filepath.Join(homeDir, ".kube", "kuberc"), true
}

func (pref kubercPreference) interferes(args parsedArgs) bool {
	if !strings.HasPrefix(pref.APIVersion, "kubectl.config.k8s.io/") || pref.Kind != "Preference" {
		return true
	}
	for _, alias := range pref.Aliases {
		if slices.Contains(args.positionals, alias.Name) {
			return true
		}
	}
	for _, d := range slices.Concat(pref.Defaults, pref.Overrides) {
		for _, option := range slices.Concat(d.Options, d.Flags) {
			if isTargetingOption(option.Name) {
				return true
			}
		}
	}
	return false
}

func isTargetingOption(name string) bool {
	name = normalizeFlagName(strings.TrimLeft(name, "-"))
	if long, ok := shortFlags[name]; ok {
		name = long
	}
	return kubercTargetingOptions[name]
}
